//w katalogu /proto
//grpc_tools_node_protoc --js_out=import_style=commonjs,binary:../gra_nodejs/src --grpc_out=generate_package_definition:../gra_nodejs/src gra.proto

const SERWER_URL = `localhost:50051`;
const grpc = require(`@grpc/grpc-js`);
const graMessages = require(`./gra_pb.js`);
const graServices = require(`./gra_grpc_pb.js`);
const readline = require('readline');

const KARTY_KTORE_WYMAGAJA_KOLORU =
{
  'L1': true,
  'L2': true,
  'A1': true,
  'A1B': true
};

const rlInterface = readline.createInterface({ input: process.stdin, output: process.stdout });
const rlIterator = rlInterface[Symbol.asyncIterator]();

(async () => {
  await main();
})();

async function main() {

  console.log('Podaj graId, jak nowa to po prostu ENTER.');
  let graIdInput = await rlIterator.next();
  let graId = graIdInput.value;

  console.log('Jak masz na imię?');
  let nazwaGraczaInput = await rlIterator.next();
  let nazwaGracza = nazwaGraczaInput.value;

  const packageDefinition = grpc.loadPackageDefinition(graServices);
  const meta = new grpc.Metadata();
  const klient = new packageDefinition.Gra(SERWER_URL, grpc.credentials.createInsecure())

  let stanGry;
  let karta;

  console.log(''); //new line

  try {

    if (!graId) {
      graId = await nowyMecz(klient, 2);
    }

    console.log(`Twój Id gry to: ${graId}`);
    console.log(`Twoje imię to: ${nazwaGracza}`);

    console.log(''); //new line
    console.log('GLHF GLHF GLHF GLHF', '\n');

    stanGry = await dolaczDoGry(klient, graId, nazwaGracza);

    while (true) {

      wyplujStanGry(stanGry);

      if (stanGry.getCzykoniec()) {
        console.log('GG', '\n');
        break;
      }

      let idKarty, idKoloru;

      // * ZAMIAST WCZYTYWANIA DO YOUR MAGIC, TWOJE KARTY, PLANSZE MACIE W STANIE GRY !!!
      // * JAK DOBRAC SIE DO STANU GRY MACIE W METODCE WYPLUJ STAN GRY
      // * ------------------------------------------------------------------------------
      
      console.log('Którą kartą chcesz zagrać?');
      let kodKartyInput = await rlIterator.next();
      const kodKarty = kodKartyInput.value;
      idKarty = konwertujKodKartyNaIdka(kodKarty);

      console.log(''); //new line

      if (KARTY_KTORE_WYMAGAJA_KOLORU[kodKarty]) {
        console.log('Podaj kolor żółwia: RED, GREEN, BLUE, YELLOW, PURPLE');
        let kolorZolwiaInput = await rlIterator.next();
        const kolorZolwia = kolorZolwiaInput.value;
        idKoloru = konwertujKolorZolwiaNaIdka(kolorZolwia)
        console.log(''); //new line
      }

      // * ------------------------------------------------------------------------------
      // * ------------------------------------------------------------------------------
      // * ------------------------------------------------------------------------------

      stanGry = await mojRuch(klient, stanGry, idKarty, idKoloru);
    }

  } catch (ex) {
    console.log(ex);
    console.log('To by było na tyle :(', '\n', '\n');
    return;
  }
}

async function nowyMecz(klient, liczbaGraczy) {
  konfiguracjaGry = new graMessages.KonfiguracjaGry();
  konfiguracjaGry.setLiczbagraczy(liczbaGraczy);

  return nowyMeczPromise = new Promise((resolve, reject) => {
    klient.nowyMecz(konfiguracjaGry, (err, res) => {
      if (err) {
        reject(err);
        return;
      }

      const graId = res.getGraid();
      if (!graId) {
        reject('Serwer nie wydał gry Id :(');
      }
      else {
        resolve(graId);
      }
    });
  });
}

async function dolaczDoGry(klient, graId, nazwaGracza) {
  dolaczanie = new graMessages.Dolaczanie();
  dolaczanie.setGraid(graId);
  dolaczanie.setNazwagracza(nazwaGracza);

  return dolaczDoGryPromise = new Promise((resolve, reject) => {
    klient.dolaczDoGry(dolaczanie, (err, res) => {
      if (err) {
        reject(err);
        return;
      }

      const stanGry = res;
      if (!stanGry) {
        reject('Serwer nie wydał stanu gry :(');
      }
      else {
        resolve(stanGry);
      }
    });
  });
}

async function mojRuch(klient, stanGry, idKarty, idKoloru) {
  ruchGracza = new graMessages.RuchGracza();
  ruchGracza.setGraid(stanGry.getGraid());
  ruchGracza.setGraczid(stanGry.getGraczid());
  ruchGracza.setZagranakarta(idKarty);
  ruchGracza.setKolorwybrany(idKoloru);

  return ruchGraczaPromise = new Promise((resolve, reject) => {
    klient.mojRuch(ruchGracza, (err, res) => {
      if (err) {
        reject(err);
        return;
      }

      const stanGry = res;
      if (!stanGry) {
        reject('Serwer nie wydał stanu gry :(');
      }
      else {
        resolve(stanGry);
      }
    });
  });
}

function wyplujStanGry(stanGry) {
  if (stanGry.getCzykoniec()) {
    console.log(`Koniec gry, wygrał gracz nr: ${stanGry.getKtowygral()}`)
  } else {
    const kolorZolwia = dajMiKolorZolwia(stanGry.getTwojkolor());
    const kartyZKodem = stanGry
      .getTwojekartyList()
      .map(k => dajMiKodKarty(k))
      .join(', ');

    console.log(`Kolor żółwia: ${kolorZolwia}`)
    console.log(`Twoje karty: ${kartyZKodem}`)

    const plansza = stanGry.getPlanszaList();
    let sytuacjaNaPlanszy = 'Plansza: ';

    plansza.forEach((pole) => {
      const zolwieNaPolu = pole.getZolwieList();
      const zolwieNaPoluPoPrzecinku = zolwieNaPolu
        .map(z => dajMiKolorZolwia(z))
        .join(', ')
        .toUpperCase();

      sytuacjaNaPlanszy += `[${zolwieNaPoluPoPrzecinku}] `;
    });

    console.log(sytuacjaNaPlanszy, '\n');
  }
}

function dajMiKolorZolwia(idKoloru) {
  return Object.keys(graMessages.KolorZolwia).find(key => graMessages.KolorZolwia[key] === idKoloru);
}

function konwertujKolorZolwiaNaIdka(kolorZolwia) {
  return graMessages.KolorZolwia[kolorZolwia];
}

function dajMiKodKarty(idKoloru) {
  return Object.keys(graMessages.Karta).find(key => graMessages.Karta[key] === idKoloru);
}

function konwertujKodKartyNaIdka(kodKarty) {
  return graMessages.Karta[kodKarty];
}
