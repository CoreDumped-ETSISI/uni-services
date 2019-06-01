<?hh // strict

use namespace HH\Lib\Regex;

/**
 * Extracts an url to the guia from the HTML code of a grado page
 */
class GuiaUrlExtractor
{
    public function __construct(
        private string $html
    )
    {
    }

    public function extract(): ?string
    {
        $pattern = re"/a href=\"(.+?)\">curso actual/";

        $m = Regex\first_match($this->html, $pattern);

        if ($m === null || count($m) < 2)
        {
            return null;
        }

        return $m[1];
    }

}
