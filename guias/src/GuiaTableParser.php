<?hh // strict

use namespace HH\Lib\Regex;
use namespace HH\Lib\C;
use namespace HH\Lib\Vec;
use namespace HH\Lib\Str;

/**
 * Parses the tables of the guia list page.
 */
class GuiaTableParser
{
    private string $html;
    
    public function __construct(string $html)
    {
        $this->html = $html;
    }

    private function createGuia(
        string $code,
        string $name,
        string $url,
        string $type,
        string $ects,
        string $semester
    ): GuiaDocente
    {
        return new GuiaDocente(
            $this->clean($code),
            $this->clean($name),
            "https://www.etsisi.upm.es" . $url,
            $this->clean($type),
            $this->clean($ects),
            $this->clean($semester)
        );
    }

    private function clean(string $str): string
    {
        $str = html_entity_decode($str);
        $str = Str\trim($str);

        // Remove tags
        $str = Regex\replace($str, re"/<\/?(?:.+?)>/", "");

        return $str;
    }

    public function parse(): vec<GuiaDocente>
    {
        $pattern = re"/<tr(?:.*?)>(.+?)<\/tr>/s";

        $matches = Regex\every_match($this->html, $pattern);

        $guias = vec<GuiaDocente>[];

        foreach ($matches as $i => $m) {
            $tr = $m[1];
            $dataP = re"/<td(?:.*?)>(.+?)<\/td>/s";

            $datas = Regex\every_match($tr, $dataP);

            if (C\count($datas) != 6) {
                continue;
            }

            $url = Regex\first_match($datas[2][1], re"/href=\"(.+?)\"/");

            if ($url === null || count($url) < 2 || $url[1] === null) {
                continue;
            }

            $guias[] = $this->createGuia(
                $datas[0][1],
                $datas[1][1],
                $url[1],
                $datas[3][1],
                $datas[4][1],
                $datas[5][1],
            );
        }

        return $guias;
    }
}
