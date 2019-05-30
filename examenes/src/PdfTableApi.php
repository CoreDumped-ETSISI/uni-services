<?hh // strict

/**
 * Submits a pdf to the microservice and parses the results.
 */
class PdfTableApi
{
    private string $url;
    private ?string $pdf;

    public function __construct(string $url)
    {
        $this->url = $url;

        $ex = new ExamPdfExtractor();
        $this->pdf = $ex->getExamUrl();
    }

    public function call(): ?vec<Table>
    {
        if ($this->pdf === null) {
            return null;
        }

        $opts = dict[
            'pdf' => $this->pdf,
            'settings' => dict[
                'pages' => '1-end' // TODO: Don't blow up the server.
                // 'pages' => '1'
            ]
        ];

        $json = json_encode($opts);
        
        $ch = curl_init($this->url);
        
        //Tell cURL that we want to send a POST request.
        curl_setopt($ch, CURLOPT_POST, 1);
        //Attach our encoded JSON string to the POST fields.
        curl_setopt($ch, CURLOPT_POSTFIELDS, $json);
        //Set the content type to application/json
        curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-Type: application/json')); 
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 300); 
        curl_setopt($ch, CURLOPT_TIMEOUT, 300); //timeout in seconds
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        
        //Execute the request
        $resp = curl_exec($ch);
        curl_close($ch);

        if (!$resp || $resp === false) {
            return null;
        }

        // Decode the response

        $tables = vec<Table>[];
        $j = json_decode($resp, true);

        foreach ($j as $i => $table) {
            // Table data is in $table['data']
            $rows = vec<vec<string>>[];
            $data = $table['data'];

            foreach ($data as $r => $row) {
                $cols = vec<string>[];

                foreach ($row as $c => $col) {
                    $cols[] = $col;
                }

                $rows[] = $cols;
            }

            $tables[] = $rows;
        }

        return $tables;
    }
}
