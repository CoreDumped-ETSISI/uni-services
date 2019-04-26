using System;
using System.Linq;
using github.Api;
using System.Threading.Tasks;

namespace github
{
    class Program
    {
        static async Task Main(string[] args)
        {
            var repos = await GitHub.GetOrganizationRepositories("CoreDumped-ETSISI");

            foreach (var repo in repos)
            {
                var langs = await GitHub.GetRepositoryLanguages(repo.FullName);

                if (langs.Count == 0) continue;

                Console.WriteLine("Repo [{2}] {0}: {1}", repo.Name, langs
                    .Select((pair) => $"{pair.Key} => {pair.Value}")
                    .Aggregate((l, r) => l + ", " + r),
                    repo.Id
                );
            }
        }
    }
}
