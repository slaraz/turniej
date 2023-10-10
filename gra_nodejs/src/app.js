//w katalogu /proto
//grpc_tools_node_protoc --js_out=import_style=commonjs,binary:../gra_nodejs/src --grpc_out=generate_package_definition:../gra_nodejs/src gra.proto

const SERWER_URL = `localhost:50051`;
const grpc = require(`@grpc/grpc-js`);
const graMessages = require(`./gra_pb.js`);
const graServices = require(`./gra_grpc_pb.js`);
const readline = require('readline');

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

  try {

    if (!graId) {
      graId = await nowyMecz(klient, 2);
    }

    console.log(`Twój Id gry to: ${graId}`);
    console.log(`Twoje imię to: ${nazwaGracza}`);
    console.log(`GLHF GLHF GLHF GLHF`);

    stanGry = await dolaczDoGry(klient, graId, nazwaGracza);

    while (true) {

      wyplujStanGry(stanGry);
      
      console.log('Którą kartą chcesz zagrać?');
      let kartaInput = await rlIterator.next();
      let karta = kartaInput.value;

      await ruchGracza(klient, stanGry, { numer: karta, kolor: 1 });
    }

  } catch (ex) {
    console.log(ex);
    console.log('To by było na tyle :(')
    return;
  }

  console.log(graId);
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

async function ruchGracza(klient, stanGry, karta) {
  ruchGracza = new graMessages.RuchGracza();
  ruchGracza.setGraid(stanGry.getGraid());
  ruchGracza.setGraczid(stanGry.getGraczid());
  ruchGracza.setZagranakarta(karta.numer);
  ruchGracza.setKolorwybrany(karta.kolor);

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
    console.log(`Twój kolor: ${stanGry.getTwojkolor()}`)
    console.log(`Twoje karty: ${stanGry.getTwojekartyList().join(', ')}`)
    // fmt.Printf("Twój kolor: %v, Pola:", stanGry.TwojKolor)
    // for _, pole := range stanGry.Plansza {
    //   fmt.Printf(" %v", pole.Zolwie)
    // }
    // fmt.Printf(", Twoje karty: %v\n", stanGry.TwojeKarty)
  }
}


