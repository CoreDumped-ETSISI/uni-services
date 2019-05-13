<?hh

require_once(__DIR__ . '/../vendor/autoload.hack');

<<__EntryPoint>>
async function main(): Awaitable<noreturn> {
    \Facebook\AutoloadMap\initialize();

    $hc = new GhettoCache(
        function() {
            $store = new ExamStore();
            $exams = $store->getExams();
            return json_encode($exams);
        }
    );

    header('Content-type:application/json;charset=utf-8');
    echo $hc->get();
    
    exit(0);
}