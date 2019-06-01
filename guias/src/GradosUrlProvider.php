<?hh // strict

/**
 * Provides URLs of different grados from which you run the guia URL extractor.
 */
class GradosUrlProvider
{
    private dict<string, string> $urls;

    public function __construct()
    {
        $this->urls = dict<string, string>[
            "Computadores" => "https://www.etsisi.upm.es/estudios/grados/61ci/organizacion-docente", // Compu
            "Software" => "https://www.etsisi.upm.es/estudios/grados/61iw/organizacion-docente", // Software
            "SI" => "https://www.etsisi.upm.es/estudios/grados/61si/organizacion-docente", // SI
            "TI" => "https://www.etsisi.upm.es/estudios/grados/61ti/organizacion-docente" // TI
        ];
    }

    public function getAllGrados(): dict<string, string>
    {
        return $this->urls;
    }
}
