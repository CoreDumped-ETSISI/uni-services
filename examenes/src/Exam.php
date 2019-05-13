<?hh // strict

/**
 * Represents an exam
 */
class Exam
{

    public function __construct(
        public string $name,
        public string $name_raw,
        public vec<string> $aulas,
        public string $timeslot,
        public string $day,
        public string $date,
        public vec<string> $tags
    )
    {}
}
