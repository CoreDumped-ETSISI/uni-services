<?hh // strict

use namespace HH\Lib\Regex;

/**
 * Extracts the PDF url from the front page.
 */
class ExamPdfExtractor
{
    public function __construct()
    {
        
    }

    public function getExamUrl(): ?string
    {
        $html = $this->downloadHomepage();

        if ($html == null) return null;

        //                               le unicode face
        $p = re"/href=\"(.*?)\"\s+title=\"E(?:.+?)s grados finales\"/";
        $m = Regex\first_match($html, $p);

        if ($m == null || count($m) < 2) {
            return null;
        }

        return $m[1];
    }

    private function downloadHomepage(): ?string
    {
        $url = "https://www.etsisi.upm.es/";

        $ch = curl_init($url);
        
        //Set the content type to application/json
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10); 
        curl_setopt($ch, CURLOPT_TIMEOUT, 10); //timeout in seconds
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        
        //Execute the request
        $resp = curl_exec($ch);
        curl_close($ch);

        if (!$resp || $resp === false) {
            return null;
        }

        return $resp;
    }
}
