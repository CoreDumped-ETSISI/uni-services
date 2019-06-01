<?hh

require_once(__DIR__ . '/../vendor/autoload.hack');

<<__EntryPoint>>
async function main(): Awaitable<noreturn> {
    \Facebook\AutoloadMap\initialize();

    $hc = new GhettoCache(
        function() {
            $ex = new GuiaExtractor();
            return $ex->json();
        }
    );

    header('Content-type:application/json;charset=utf-8');
    echo $hc->get();
    
    exit(0);
}