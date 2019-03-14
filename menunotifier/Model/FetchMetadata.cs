using System;
using System.Collections.Generic;

namespace menunotifier.Model
{
    public class FetchMetadata
    {
        public List<Menu> Menu { get; set; }
        public string Mes { get; set; }
        public int Desde { get; set; }
        public int Hasta { get; set; }
    }
}