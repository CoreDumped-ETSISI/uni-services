<?hh // strict

/**
 * Extracts the PDF url from the front page.
 */
class ExamPdfExtractor
{
    public function __construct()
    {
        
    }

    public function getExamUrl(): string
    {
        // TODO: actually extract the url.
        return 'http://www.etsisi.upm.es/sites/default/files/curso_2018_19/Grado_Planificacion/examenes_finales.pdf';
    }
}
