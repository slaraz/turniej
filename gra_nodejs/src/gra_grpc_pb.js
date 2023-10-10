// GENERATED CODE -- DO NOT EDIT!

'use strict';
var gra_pb = require('./gra_pb.js');

function serialize_Dolaczanie(arg) {
  if (!(arg instanceof gra_pb.Dolaczanie)) {
    throw new Error('Expected argument of type Dolaczanie');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Dolaczanie(buffer_arg) {
  return gra_pb.Dolaczanie.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_KonfiguracjaGry(arg) {
  if (!(arg instanceof gra_pb.KonfiguracjaGry)) {
    throw new Error('Expected argument of type KonfiguracjaGry');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_KonfiguracjaGry(buffer_arg) {
  return gra_pb.KonfiguracjaGry.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_NowaGraInfo(arg) {
  if (!(arg instanceof gra_pb.NowaGraInfo)) {
    throw new Error('Expected argument of type NowaGraInfo');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_NowaGraInfo(buffer_arg) {
  return gra_pb.NowaGraInfo.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_RuchGracza(arg) {
  if (!(arg instanceof gra_pb.RuchGracza)) {
    throw new Error('Expected argument of type RuchGracza');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_RuchGracza(buffer_arg) {
  return gra_pb.RuchGracza.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_StanGry(arg) {
  if (!(arg instanceof gra_pb.StanGry)) {
    throw new Error('Expected argument of type StanGry');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_StanGry(buffer_arg) {
  return gra_pb.StanGry.deserializeBinary(new Uint8Array(buffer_arg));
}


var GraService = exports['Gra'] = {
  // tworzy mecz
// zwraca ID gry
nowyMecz: {
    path: '/Gra/NowyMecz',
    requestStream: false,
    responseStream: false,
    requestType: gra_pb.KonfiguracjaGry,
    responseType: gra_pb.NowaGraInfo,
    requestSerialize: serialize_KonfiguracjaGry,
    requestDeserialize: deserialize_KonfiguracjaGry,
    responseSerialize: serialize_NowaGraInfo,
    responseDeserialize: deserialize_NowaGraInfo,
  },
  // dołącza do nowej gry
dolaczDoGry: {
    path: '/Gra/DolaczDoGry',
    requestStream: false,
    responseStream: false,
    requestType: gra_pb.Dolaczanie,
    responseType: gra_pb.StanGry,
    requestSerialize: serialize_Dolaczanie,
    requestDeserialize: deserialize_Dolaczanie,
    responseSerialize: serialize_StanGry,
    responseDeserialize: deserialize_StanGry,
  },
  // ruch gracza klienta
mojRuch: {
    path: '/Gra/MojRuch',
    requestStream: false,
    responseStream: false,
    requestType: gra_pb.RuchGracza,
    responseType: gra_pb.StanGry,
    requestSerialize: serialize_RuchGracza,
    requestDeserialize: deserialize_RuchGracza,
    responseSerialize: serialize_StanGry,
    responseDeserialize: deserialize_StanGry,
  },
};

