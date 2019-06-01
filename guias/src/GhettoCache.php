<?hh // strict

/**
 * Caches a string on the disk.
 */
class GhettoCache
{
    public function __construct(
        private (function():string) $callback
    )
    {}

    private function exists(): bool
    {
        return file_exists(".cache");
    }

    public function get(): string
    {
        if ($this->exists()) {
            return file_get_contents(".cache");
        }

        $cb = $this->callback;
        
        $str = $cb();

        file_put_contents(".cache", $str);

        return $str;
    }
}
