<?hh // strict

/**
 * Simple HTTP Client to download data
 */
class HttpClient
{
    public function __construct()
    {
                
    }

    public function get(string $url): ?string
    {
        $ch = curl_init($url);
        
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
