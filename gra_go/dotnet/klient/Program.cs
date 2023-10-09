using Grpc.Net.Client;
using GraKlient;
using CommandLine;
using Grpc.Core;
using System.Linq.Expressions;
using System.Diagnostics;

internal class Program
{
    private static void Main(string[] args)
    {
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
                            LiczbaGraczy = 2
                        });
                        Console.WriteLine($"Nowa gra: {graInfo.GraID} (Opis: {graInfo.Opis})");
                        graId = graInfo.GraID;
                    }

                    Console.WriteLine($"Gracz {o.Nazwa} dołącza do gry {graId}{Environment.NewLine}Czekam na stan gry...");
                    var stanGry = client.DolaczDoGry(new Dolaczanie()
                    {
                        GraID = graId,
                        Wizytowka = new WizytowkaGracza() { Nazwa = o.Nazwa }
                    });
                    Console.WriteLine($"Stan gry: plansza: {stanGry.SytuacjaNaPlanszy}, karty: {stanGry.TwojeKarty}");

                    while (!CzyKoniec(stanGry))
                    {
                        Console.WriteLine("Wybierz kartę do zagrania:");
                        Console.Write(">");
                        var karta = Console.ReadLine();

                        Console.WriteLine($"Gracz {stanGry.GraID}-{stanGry.GraczID} zagrywa kartę: {karta}");
                        Console.WriteLine("Czekam na stan gry...");

                        stanGry = client.MojRuch(new RuchGracza() {
                            GraID = stanGry.GraID,
                            GraczID = stanGry.GraczID,
                            ZagranaKarta = karta,
                        });
                    }

                });
            Console.WriteLine("Koniec");
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }
    }

    private static bool CzyKoniec(StanGry stanGry) => stanGry.SytuacjaNaPlanszy == "KONIEC";

    private static T SafeRun<T>(Func<T> action)
    {
        try
        {
            return action();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Wystąpiła błąd: {ex.Message}");
            throw;
        }
    }
}