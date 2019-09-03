<?hh // strict

type Table = vec<vec<string>>;

use namespace HH\Lib\C;
use namespace HH\Lib\Vec;
use namespace HH\Lib\Str;
use namespace HH\Lib\Regex;
use namespace HH\Lib\Set;

class ParsingException extends Exception
{

}

/**
 * Parses the PDF tables returned from the extractor microservice
 */
class ExamPdfParser
{
    private vec<Table> $tables;
    
    public function __construct(vec<Table> $tables)
    {
        $this->tables = $tables;
    }
    
    private function sanitize(string $str): string
    {
        return Str\replace($str, "\n", " ");
    }

    private function parseTimeslot(string $slot): string
    {
        $slots = dict[
            'M1' => '9:00 a 12:00',
            'M2' => '12:00 a 15:00',
            'T1' => '15:00 a 18:00',
            'T2' => '18:00 a 21:00',
        ];

        if (!C\contains_key($slots, $slot)) {
            throw new ParsingException();
        }

        return $slots[$slot];
    }

    private function parseTime(string $datetime): string
    {
        $day = Regex\first_match($datetime, re"/\d+/")[0];
        $month = Regex\first_match($datetime, re"/[A-Za-z]+$/")[0];

        if ($month === null) {
            throw new ParsingException();
        }

        return $day . " de " . $month;
    }

    private function getDate(string $datetime, string $year): string
    {
        $datetime = Str\trim($datetime);
        $day = Regex\first_match($datetime, re"/\d+/")[0];
        $month = Regex\first_match($datetime, re"/[A-Za-z]+$/")[0];

        if ($month === null) {
            throw new ParsingException();
        }

        $month = Str\lowercase($month);

        $months = dict[
            'enero' => 1,
            'febrero' => 2,
            'marzo' => 3,
            'abril' => 4,
            'mayo' => 5,
            'junio' => 6,
            'julio' => 7,
            'agosto' => 8,
            'septiembre' => 9,
            'octubre' => 10,
            'noviembre' => 11,
            'diciembre' => 12,
        ];

        if (!C\contains_key($months, $month)) {
            throw new ParsingException();
        }

        return $day . '/' . $months[$month] . '/' . $year;
    }

    private function trimTags(vec<string> $tags): vec<string>
    {
        return $tags
            |> Vec\filter($$, $t ==> $t !== "y")
            |> Vec\map($$, $t ==> Str\trim($t, ",()"));
    }

    private function getTags(string $name): vec<string>
    {
        $raw_tags = vec<string>[];

        $grados = Regex\every_match($name, re"/\([A-Za-z0-9,\s]+\)/");
        
        foreach ($grados as $i => $match) {
            $grado = $match[0];
            $trimmed = $this->trimTags(Str\split($grado, " "));
            $raw_tags = Vec\concat($raw_tags, $trimmed);
        }

        // NOTICE: Plan 2009 is not supported.
        $tag_activators = dict[
            "general" => function(string $t): bool {
                return Regex\matches($t, re"/G\d+/") ||
                    Regex\matches($t, re"/GOPT$/");
            },
            "software" => function(string $t): bool {
                return Regex\matches($t, re"/GS\d+/") ||
                        Regex\matches($t, re"/G[A-Z]*S(?!I)[A-Z]*OPT$/");
            },
            "compu" => function(string $t): bool {
                return Regex\matches($t, re"/GC\d+/") || 
                    Regex\matches($t, re"/G[A-Z]*C[A-Z]*OPT$/");
            },
            "si" => function(string $t): bool {
                return Regex\matches($t, re"/GSI\d+/") || 
                    Regex\matches($t, re"/G[A-Z]*SI[A-Z]*OPT$/");
            },
            "ti" => function(string $t): bool {
                return Regex\matches($t, re"/GTI\d+/") || 
                    Regex\matches($t, re"/G[A-Z]*TI[A-Z]*OPT$/");
            },
            "optativa" => function(string $t): bool {
                return Regex\matches($t, re"/G[A-Z]*OPT$/");
            },
        ];

        $tags = Set{};

        foreach ($raw_tags as $i => $tag) {
            foreach ($tag_activators as $tagname => $f) {
                if ($f($tag)) {
                    $tags[] = $tagname;

                    $year = Regex\first_match($tag, re"/\d+/");

                    if ($year !== null) {
                        $tags[] = "year_" . $year[0];
                    }
                }
            }

        }


        return vec<string>($tags);
    }

    private function getPureName(string $name): string
    {
        return $name
            |> Regex\replace($$, re"/\([A-Za-z0-9,\s]+\)/", "")
            |> Str\trim($$)
            |> $this->sanitize($$);
    }

    private function createExam(string $name, string $block, string $slot, string $day, string $year): Exam
    {
        return new Exam(
            $this->getPureName($name),
            $this->sanitize($name),
            Str\split($block, ",")
                |> Vec\map($$, $b ==> $this->sanitize($b))
                |> Vec\map($$, $b ==> Str\trim($b)),
            $this->parseTimeslot($slot),
            $this->parseTime($day),
            $this->getDate($day, $year),
            $this->getTags($name)
        );
    }

    private function getYear(string $header): string {
        $y = Regex\first_match($header, re"/\d+$/")[0];
        return $y;
    }

    public function parse(): vec<Exam>
    {
        $exams = vec<Exam>[];
        $year = "";

        foreach ($this->tables as $i => $table) {
            // First of all, filter out the calendar tables.
            if (C\count($table) > 1 && C\count($table[1]) > 0 && $table[1][0] === "L")
            {
                $year = $this->getYear($table[0][0]);
                continue;
            }
            
            $day = "";
            
            foreach ($table as $r => $row) {
                if (C\count($row) < 6) {
                    // wtf is this even
                    continue;
                }
                // Skip table headers
                if ($row[1] === "TURNO") {
                    continue;
                }

                if ($row[0] !== "") {
                    $day = $this->sanitize($row[0]);
                }

                $slot = $row[1];

                $lex = $row[2];
                $lblock = $row[3];

                $rex = $row[4];
                $rblock = $row[5];

                if ($lex) {
                    try {
                        $exams[] = $this->createExam(
                            $lex,
                            $lblock,
                            $slot,
                            $day,
                            $year
                        );
                    } catch(ParsingException $ex) {
                        // Continue
                    }
                }

                if ($rex) {
                    try {
                        $exams[] = $this->createExam(
                            $rex,
                            $rblock,
                            $slot,
                            $day,
                            $year
                        );
                    } catch(ParsingException $ex) {
                        // Continue
                    }
                    
                }
            }
        }

        return $exams;
    }
}
