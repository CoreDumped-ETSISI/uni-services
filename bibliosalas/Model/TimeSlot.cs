namespace bibliosalas.Model
{
    public class TimeSlot
    {
        public const int HourStart = 9;
        public const int HourEnd = 20;

        public int Start { get; set; }
        public int End { get; set; }

        public override string ToString()
        {
            return string.Format("{0}:{1:D2} - {2}:{3:D2}", Start / 100, Start % 100, End / 100, End % 100);
        }

        public static int IndexToTime(int index)
        {
            int h = TimeSlot.HourStart + (index / 2) * 1;
            int m = (index % 2) * 30;

            return h * 100 + m;
        }
    }
}