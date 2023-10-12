using Grpc.Net.Client;
using GraKlient;
using CommandLine;

try
{
    Console.WriteLine("Start");
    Parser.Default.ParseArguments<Options>(args)
        .WithParsed(o =>
        {
            using var channel = GrpcChannel.ForAddress(o.Serwer);
            var client = new Gra.GraClient(channel);

            var graId = o.GraId;
            if (o.CzyNowa)
            {
                var graInfo = client.NowyMecz(new KonfiguracjaGry()
                {
                    LiczbaGraczy = o.LiczbaGraczy
                });
                Console.WriteLine($"Nowa gra: {graInfo.GraID}");
                graId = graInfo.GraID;
            }

            Console.WriteLine($"Gracz {o.Nazwa} dołącza do gry {graId}{Environment.NewLine}Czekam na odpowiedź od serwera...");
            var stanGry = client.DolaczDoGry(new Dolaczanie()
            {
                GraID = graId,
                NazwaGracza = o.Nazwa,
            });

            var kartyDlaKtorychTrzebaPodacKolor = new List<Karta>() { Karta.L1, Karta.L2, Karta.A1, Karta.A1B };
            Karta karta;
            KolorZolwia kolor;

            while (true)
            {
                DrukujStatus(stanGry);
                if (stanGry.CzyKoniec)
                {
                    break;
                }

                karta = WczytajKarte();
                karta = UpewnijSieZeKartaPoprawna(stanGry, karta);

                if (kartyDlaKtorychTrzebaPodacKolor.Contains(karta))
                {
                    kolor = WczytajKolor();
                }
                else
                {
                    kolor = KolorZolwia.Xxx;
                }

                Console.WriteLine($"Gracz {stanGry.GraID}-{stanGry.GraczID} zagrywa kartę: {karta}");
                Console.WriteLine("Czekam na odpowiedź serwera...");

                stanGry = client.MojRuch(new RuchGracza()
                {
                    GraID = stanGry.GraID,
                    GraczID = stanGry.GraczID,
                    ZagranaKarta = karta,
                    KolorWybrany = kolor,
                });
            }

        });
    Console.WriteLine("Koniec");
}
catch (Exception ex)
{
    Console.WriteLine(ex.Message);
}


static Karta UpewnijSieZeKartaPoprawna(StanGry stanGry, Karta karta)
{
    while (!stanGry.TwojeKarty.Contains(karta))
    {
        Console.WriteLine("Ta karta nie może zostać użyta na tym etapie. Wprowadź poprawną kartę:");
        karta = WczytajKarte();
    };

    return karta;
}

static KolorZolwia WczytajKolor()
{
    Console.WriteLine("Wybierz kolor:");
    var ans = Console.ReadLine() ?? string.Empty;

    KolorZolwia kolorZolwia;
    while (!Enum.TryParse<KolorZolwia>(Capitalize(ans), out kolorZolwia))
    {
        Console.Write("Niepoprawny kolor żółwia. Proszę wprowadź kolor ponownie:");
        ans = Console.ReadLine();
    }

    return kolorZolwia;
}

static Karta WczytajKarte()
{
    Console.Write("Wybierz kartę do zagrania: ");
    var ans = Console.ReadLine() ?? string.Empty;

    Karta karta;
    while (!Enum.TryParse<Karta>(Capitalize(ans), out karta))
    {
        Console.Write("Karta niepoprawna. Proszę wprowadź kartę ponownie:");
        ans = Console.ReadLine();
    }

    return karta;
}

static void DrukujStatus(StanGry stanGry)
{
    if (stanGry.CzyKoniec)
    {
        Console.WriteLine($"Koniec gry, wygrał gracz nr {stanGry.KtoWygral}");
    }
    else
    {
        Console.WriteLine($"Twój kolor: {stanGry.TwojKolor}, Pola: {string.Join(", ", stanGry.Plansza.Select(x => x.Zolwie))}, Twoje karty: {stanGry.TwojeKarty} ");
    }
}

static string Capitalize(string str)
{
    if (str is null)
        return string.Empty;

    if (str.Length > 1)
        return str.Substring(0, 1).ToUpper() + str.Substring(1).ToLower();

    return str.ToUpper(); ;
}
