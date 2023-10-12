using CommandLine.Text;
using CommandLine;

public class Options
{
    [Option("nazwa", Required = false, HelpText = "nazwa gracza")]
    public string Nazwa { get; set; } = "Ziutek .Net";

    [Option('a', "addr", Required = false, HelpText = "wskazuje nazwę serwera do podłączenia")]
    public string Serwer { get; set; } = "http://localhost:50051";

    [Option('n', "nowa", Required = true, HelpText = "rozpoczyna nową grę", SetName = "tryb1")]
    public bool CzyNowa { get; set; }
    
    [Option('g', "gra", Required = true, HelpText = "dołącza do gry o podanej nazwie", SetName = "tryb2")]
    public string GraId { get; set; }

    [Option('l', "lg", Required = false, HelpText = "określa liczbę graczy")]
    public int LiczbaGraczy { get; set; } = 2;
}