<?hh // strict

/**
 * Manages the storage of exams
 */
class ExamStore
{
    private ?vec<Exam> $exams;

    public function __construct()
    {
        $this->exams = null;
    }

    public function getExams(): vec<Exam>
    {
        if ($this->exams !== null) {
            return $this->exams;
        }

        $endpoint = getenv("PDFTABLE_SERVER");
        $api = new PdfTableApi($endpoint);
        $tables = $api->call();

        if ($tables === null) {
            return vec<Exam>[];
        }

        $parser = new ExamPdfParser($tables);

        return $parser->parse();
    }
}
