<?hh // strict

/**
 * Guia for a class
 */
class GuiaDocente
{
    public function __construct(
        public string $code,
        public string $name,
        public string $url,
        public string $type,
        public string $ects,
        public string $semestre
    )
    {
        
    }
}
