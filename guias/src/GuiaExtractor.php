<?hh // strict

/**
 * Extracts all the guias and optionally puts them in a json.
 */
class GuiaExtractor
{
    public function getGuias(): dict<string, vec<GuiaDocente>>
    {
        $wc = new HttpClient();
        $d = dict<string, vec<GuiaDocente>>[];

        $p = new GradosUrlProvider();

        $grados = $p->getAllGrados();

        foreach ($grados as $grado => $url) {
            $html = $wc->get($url);

            if ($html === null) {
                continue;
            }

            $gext = new GuiaUrlExtractor($html);
            $tableurl = $gext->extract();

            if ($tableurl === null) {
                continue;
            }

            $tablehtml = $wc->get($tableurl);
            if ($tablehtml is nonnull) {
                $parser = new GuiaTableParser($tablehtml);
                $guias = $parser->parse();

                $d[$grado] = $guias;
            }
        }

        return $d;
    }

    public function json(): string
    {
        $d = $this->getGuias();

        return json_encode($d);
    }
}
