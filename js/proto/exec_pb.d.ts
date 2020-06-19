// package: exec
// file: exec.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as github_com_gogo_protobuf_gogoproto_gogo_pb from "./github.com/gogo/protobuf/gogoproto/gogo_pb";
import * as github_com_tendermint_tendermint_abci_types_types_pb from "./github.com/tendermint/tendermint/abci/types/types_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as errors_pb from "./errors_pb";
import * as names_pb from "./names_pb";
import * as txs_pb from "./txs_pb";
import * as permission_pb from "./permission_pb";
import * as spec_pb from "./spec_pb";

export class StreamEvents extends jspb.Message { 
    clearStreameventsList(): void;
    getStreameventsList(): Array<StreamEvent>;
    setStreameventsList(value: Array<StreamEvent>): void;
    addStreamevents(value?: StreamEvent, index?: number): StreamEvent;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StreamEvents.AsObject;
    static toObject(includeInstance: boolean, msg: StreamEvents): StreamEvents.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StreamEvents, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StreamEvents;
    static deserializeBinaryFromReader(message: StreamEvents, reader: jspb.BinaryReader): StreamEvents;
}

export namespace StreamEvents {
    export type AsObject = {
        streameventsList: Array<StreamEvent.AsObject>,
    }
}

export class StreamEvent extends jspb.Message { 

    hasBeginblock(): boolean;
    clearBeginblock(): void;
    getBeginblock(): BeginBlock | undefined;
    setBeginblock(value?: BeginBlock): void;


    hasBegintx(): boolean;
    clearBegintx(): void;
    getBegintx(): BeginTx | undefined;
    setBegintx(value?: BeginTx): void;


    hasEnvelope(): boolean;
    clearEnvelope(): void;
    getEnvelope(): txs_pb.Envelope | undefined;
    setEnvelope(value?: txs_pb.Envelope): void;


    hasEvent(): boolean;
    clearEvent(): void;
    getEvent(): Event | undefined;
    setEvent(value?: Event): void;


    hasEndtx(): boolean;
    clearEndtx(): void;
    getEndtx(): EndTx | undefined;
    setEndtx(value?: EndTx): void;


    hasEndblock(): boolean;
    clearEndblock(): void;
    getEndblock(): EndBlock | undefined;
    setEndblock(value?: EndBlock): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StreamEvent.AsObject;
    static toObject(includeInstance: boolean, msg: StreamEvent): StreamEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StreamEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StreamEvent;
    static deserializeBinaryFromReader(message: StreamEvent, reader: jspb.BinaryReader): StreamEvent;
}

export namespace StreamEvent {
    export type AsObject = {
        beginblock?: BeginBlock.AsObject,
        begintx?: BeginTx.AsObject,
        envelope?: txs_pb.Envelope.AsObject,
        event?: Event.AsObject,
        endtx?: EndTx.AsObject,
        endblock?: EndBlock.AsObject,
    }
}

export class BeginBlock extends jspb.Message { 
    getHeight(): number;
    setHeight(value: number): void;

    getNumtxs(): number;
    setNumtxs(value: number): void;

    getPredecessorheight(): number;
    setPredecessorheight(value: number): void;


    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): github_com_tendermint_tendermint_abci_types_types_pb.Header | undefined;
    setHeader(value?: github_com_tendermint_tendermint_abci_types_types_pb.Header): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BeginBlock.AsObject;
    static toObject(includeInstance: boolean, msg: BeginBlock): BeginBlock.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BeginBlock, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BeginBlock;
    static deserializeBinaryFromReader(message: BeginBlock, reader: jspb.BinaryReader): BeginBlock;
}

export namespace BeginBlock {
    export type AsObject = {
        height: number,
        numtxs: number,
        predecessorheight: number,
        header?: github_com_tendermint_tendermint_abci_types_types_pb.Header.AsObject,
    }
}

export class EndBlock extends jspb.Message { 
    getHeight(): number;
    setHeight(value: number): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EndBlock.AsObject;
    static toObject(includeInstance: boolean, msg: EndBlock): EndBlock.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EndBlock, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EndBlock;
    static deserializeBinaryFromReader(message: EndBlock, reader: jspb.BinaryReader): EndBlock;
}

export namespace EndBlock {
    export type AsObject = {
        height: number,
    }
}

export class BeginTx extends jspb.Message { 

    hasTxheader(): boolean;
    clearTxheader(): void;
    getTxheader(): TxHeader | undefined;
    setTxheader(value?: TxHeader): void;

    getNumevents(): number;
    setNumevents(value: number): void;


    hasResult(): boolean;
    clearResult(): void;
    getResult(): Result | undefined;
    setResult(value?: Result): void;


    hasException(): boolean;
    clearException(): void;
    getException(): errors_pb.Exception | undefined;
    setException(value?: errors_pb.Exception): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BeginTx.AsObject;
    static toObject(includeInstance: boolean, msg: BeginTx): BeginTx.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BeginTx, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BeginTx;
    static deserializeBinaryFromReader(message: BeginTx, reader: jspb.BinaryReader): BeginTx;
}

export namespace BeginTx {
    export type AsObject = {
        txheader?: TxHeader.AsObject,
        numevents: number,
        result?: Result.AsObject,
        exception?: errors_pb.Exception.AsObject,
    }
}

export class EndTx extends jspb.Message { 
    getTxhash(): Uint8Array | string;
    getTxhash_asU8(): Uint8Array;
    getTxhash_asB64(): string;
    setTxhash(value: Uint8Array | string): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EndTx.AsObject;
    static toObject(includeInstance: boolean, msg: EndTx): EndTx.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EndTx, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EndTx;
    static deserializeBinaryFromReader(message: EndTx, reader: jspb.BinaryReader): EndTx;
}

export namespace EndTx {
    export type AsObject = {
        txhash: Uint8Array | string,
    }
}

export class TxHeader extends jspb.Message { 
    getTxtype(): number;
    setTxtype(value: number): void;

    getTxhash(): Uint8Array | string;
    getTxhash_asU8(): Uint8Array;
    getTxhash_asB64(): string;
    setTxhash(value: Uint8Array | string): void;

    getHeight(): number;
    setHeight(value: number): void;

    getIndex(): number;
    setIndex(value: number): void;


    hasOrigin(): boolean;
    clearOrigin(): void;
    getOrigin(): Origin | undefined;
    setOrigin(value?: Origin): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TxHeader.AsObject;
    static toObject(includeInstance: boolean, msg: TxHeader): TxHeader.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TxHeader, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TxHeader;
    static deserializeBinaryFromReader(message: TxHeader, reader: jspb.BinaryReader): TxHeader;
}

export namespace TxHeader {
    export type AsObject = {
        txtype: number,
        txhash: Uint8Array | string,
        height: number,
        index: number,
        origin?: Origin.AsObject,
    }
}

export class BlockExecution extends jspb.Message { 
    getHeight(): number;
    setHeight(value: number): void;

    getPredecessorheight(): number;
    setPredecessorheight(value: number): void;


    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): github_com_tendermint_tendermint_abci_types_types_pb.Header | undefined;
    setHeader(value?: github_com_tendermint_tendermint_abci_types_types_pb.Header): void;

    clearTxexecutionsList(): void;
    getTxexecutionsList(): Array<TxExecution>;
    setTxexecutionsList(value: Array<TxExecution>): void;
    addTxexecutions(value?: TxExecution, index?: number): TxExecution;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BlockExecution.AsObject;
    static toObject(includeInstance: boolean, msg: BlockExecution): BlockExecution.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BlockExecution, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BlockExecution;
    static deserializeBinaryFromReader(message: BlockExecution, reader: jspb.BinaryReader): BlockExecution;
}

export namespace BlockExecution {
    export type AsObject = {
        height: number,
        predecessorheight: number,
        header?: github_com_tendermint_tendermint_abci_types_types_pb.Header.AsObject,
        txexecutionsList: Array<TxExecution.AsObject>,
    }
}

export class TxExecutionKey extends jspb.Message { 
    getHeight(): number;
    setHeight(value: number): void;

    getOffset(): number;
    setOffset(value: number): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TxExecutionKey.AsObject;
    static toObject(includeInstance: boolean, msg: TxExecutionKey): TxExecutionKey.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TxExecutionKey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TxExecutionKey;
    static deserializeBinaryFromReader(message: TxExecutionKey, reader: jspb.BinaryReader): TxExecutionKey;
}

export namespace TxExecutionKey {
    export type AsObject = {
        height: number,
        offset: number,
    }
}

export class TxExecution extends jspb.Message { 

    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): TxHeader | undefined;
    setHeader(value?: TxHeader): void;


    hasEnvelope(): boolean;
    clearEnvelope(): void;
    getEnvelope(): txs_pb.Envelope | undefined;
    setEnvelope(value?: txs_pb.Envelope): void;

    clearEventsList(): void;
    getEventsList(): Array<Event>;
    setEventsList(value: Array<Event>): void;
    addEvents(value?: Event, index?: number): Event;


    hasResult(): boolean;
    clearResult(): void;
    getResult(): Result | undefined;
    setResult(value?: Result): void;


    hasReceipt(): boolean;
    clearReceipt(): void;
    getReceipt(): txs_pb.Receipt | undefined;
    setReceipt(value?: txs_pb.Receipt): void;


    hasException(): boolean;
    clearException(): void;
    getException(): errors_pb.Exception | undefined;
    setException(value?: errors_pb.Exception): void;

    clearTxexecutionsList(): void;
    getTxexecutionsList(): Array<TxExecution>;
    setTxexecutionsList(value: Array<TxExecution>): void;
    addTxexecutions(value?: TxExecution, index?: number): TxExecution;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TxExecution.AsObject;
    static toObject(includeInstance: boolean, msg: TxExecution): TxExecution.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TxExecution, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TxExecution;
    static deserializeBinaryFromReader(message: TxExecution, reader: jspb.BinaryReader): TxExecution;
}

export namespace TxExecution {
    export type AsObject = {
        header?: TxHeader.AsObject,
        envelope?: txs_pb.Envelope.AsObject,
        eventsList: Array<Event.AsObject>,
        result?: Result.AsObject,
        receipt?: txs_pb.Receipt.AsObject,
        exception?: errors_pb.Exception.AsObject,
        txexecutionsList: Array<TxExecution.AsObject>,
    }
}

export class Origin extends jspb.Message { 
    getChainid(): string;
    setChainid(value: string): void;

    getHeight(): number;
    setHeight(value: number): void;

    getIndex(): number;
    setIndex(value: number): void;


    hasTime(): boolean;
    clearTime(): void;
    getTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTime(value?: google_protobuf_timestamp_pb.Timestamp): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Origin.AsObject;
    static toObject(includeInstance: boolean, msg: Origin): Origin.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Origin, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Origin;
    static deserializeBinaryFromReader(message: Origin, reader: jspb.BinaryReader): Origin;
}

export namespace Origin {
    export type AsObject = {
        chainid: string,
        height: number,
        index: number,
        time?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}

export class Header extends jspb.Message { 
    getTxtype(): number;
    setTxtype(value: number): void;

    getTxhash(): Uint8Array | string;
    getTxhash_asU8(): Uint8Array;
    getTxhash_asB64(): string;
    setTxhash(value: Uint8Array | string): void;

    getEventtype(): number;
    setEventtype(value: number): void;

    getEventid(): string;
    setEventid(value: string): void;

    getHeight(): number;
    setHeight(value: number): void;

    getIndex(): number;
    setIndex(value: number): void;


    hasException(): boolean;
    clearException(): void;
    getException(): errors_pb.Exception | undefined;
    setException(value?: errors_pb.Exception): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Header.AsObject;
    static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Header;
    static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
}

export namespace Header {
    export type AsObject = {
        txtype: number,
        txhash: Uint8Array | string,
        eventtype: number,
        eventid: string,
        height: number,
        index: number,
        exception?: errors_pb.Exception.AsObject,
    }
}

export class Event extends jspb.Message { 

    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): Header | undefined;
    setHeader(value?: Header): void;


    hasInput(): boolean;
    clearInput(): void;
    getInput(): InputEvent | undefined;
    setInput(value?: InputEvent): void;


    hasOutput(): boolean;
    clearOutput(): void;
    getOutput(): OutputEvent | undefined;
    setOutput(value?: OutputEvent): void;


    hasCall(): boolean;
    clearCall(): void;
    getCall(): CallEvent | undefined;
    setCall(value?: CallEvent): void;


    hasLog(): boolean;
    clearLog(): void;
    getLog(): LogEvent | undefined;
    setLog(value?: LogEvent): void;


    hasGovernaccount(): boolean;
    clearGovernaccount(): void;
    getGovernaccount(): GovernAccountEvent | undefined;
    setGovernaccount(value?: GovernAccountEvent): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Event.AsObject;
    static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Event;
    static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
    export type AsObject = {
        header?: Header.AsObject,
        input?: InputEvent.AsObject,
        output?: OutputEvent.AsObject,
        call?: CallEvent.AsObject,
        log?: LogEvent.AsObject,
        governaccount?: GovernAccountEvent.AsObject,
    }
}

export class Result extends jspb.Message { 
    getReturn(): Uint8Array | string;
    getReturn_asU8(): Uint8Array;
    getReturn_asB64(): string;
    setReturn(value: Uint8Array | string): void;

    getGasused(): number;
    setGasused(value: number): void;


    hasNameentry(): boolean;
    clearNameentry(): void;
    getNameentry(): names_pb.Entry | undefined;
    setNameentry(value?: names_pb.Entry): void;


    hasPermargs(): boolean;
    clearPermargs(): void;
    getPermargs(): permission_pb.PermArgs | undefined;
    setPermargs(value?: permission_pb.PermArgs): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Result.AsObject;
    static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Result;
    static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
}

export namespace Result {
    export type AsObject = {
        pb_return: Uint8Array | string,
        gasused: number,
        nameentry?: names_pb.Entry.AsObject,
        permargs?: permission_pb.PermArgs.AsObject,
    }
}

export class LogEvent extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): void;

    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): void;

    clearTopicsList(): void;
    getTopicsList(): Array<Uint8Array | string>;
    getTopicsList_asU8(): Array<Uint8Array>;
    getTopicsList_asB64(): Array<string>;
    setTopicsList(value: Array<Uint8Array | string>): void;
    addTopics(value: Uint8Array | string, index?: number): Uint8Array | string;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LogEvent.AsObject;
    static toObject(includeInstance: boolean, msg: LogEvent): LogEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LogEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LogEvent;
    static deserializeBinaryFromReader(message: LogEvent, reader: jspb.BinaryReader): LogEvent;
}

export namespace LogEvent {
    export type AsObject = {
        address: Uint8Array | string,
        data: Uint8Array | string,
        topicsList: Array<Uint8Array | string>,
    }
}

export class CallEvent extends jspb.Message { 
    getCalltype(): number;
    setCalltype(value: number): void;


    hasCalldata(): boolean;
    clearCalldata(): void;
    getCalldata(): CallData | undefined;
    setCalldata(value?: CallData): void;

    getOrigin(): Uint8Array | string;
    getOrigin_asU8(): Uint8Array;
    getOrigin_asB64(): string;
    setOrigin(value: Uint8Array | string): void;

    getStackdepth(): number;
    setStackdepth(value: number): void;

    getReturn(): Uint8Array | string;
    getReturn_asU8(): Uint8Array;
    getReturn_asB64(): string;
    setReturn(value: Uint8Array | string): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CallEvent.AsObject;
    static toObject(includeInstance: boolean, msg: CallEvent): CallEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CallEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CallEvent;
    static deserializeBinaryFromReader(message: CallEvent, reader: jspb.BinaryReader): CallEvent;
}

export namespace CallEvent {
    export type AsObject = {
        calltype: number,
        calldata?: CallData.AsObject,
        origin: Uint8Array | string,
        stackdepth: number,
        pb_return: Uint8Array | string,
    }
}

export class GovernAccountEvent extends jspb.Message { 

    hasAccountupdate(): boolean;
    clearAccountupdate(): void;
    getAccountupdate(): spec_pb.TemplateAccount | undefined;
    setAccountupdate(value?: spec_pb.TemplateAccount): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GovernAccountEvent.AsObject;
    static toObject(includeInstance: boolean, msg: GovernAccountEvent): GovernAccountEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GovernAccountEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GovernAccountEvent;
    static deserializeBinaryFromReader(message: GovernAccountEvent, reader: jspb.BinaryReader): GovernAccountEvent;
}

export namespace GovernAccountEvent {
    export type AsObject = {
        accountupdate?: spec_pb.TemplateAccount.AsObject,
    }
}

export class InputEvent extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InputEvent.AsObject;
    static toObject(includeInstance: boolean, msg: InputEvent): InputEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InputEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InputEvent;
    static deserializeBinaryFromReader(message: InputEvent, reader: jspb.BinaryReader): InputEvent;
}

export namespace InputEvent {
    export type AsObject = {
        address: Uint8Array | string,
    }
}

export class OutputEvent extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OutputEvent.AsObject;
    static toObject(includeInstance: boolean, msg: OutputEvent): OutputEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OutputEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OutputEvent;
    static deserializeBinaryFromReader(message: OutputEvent, reader: jspb.BinaryReader): OutputEvent;
}

export namespace OutputEvent {
    export type AsObject = {
        address: Uint8Array | string,
    }
}

export class CallData extends jspb.Message { 
    getCaller(): Uint8Array | string;
    getCaller_asU8(): Uint8Array;
    getCaller_asB64(): string;
    setCaller(value: Uint8Array | string): void;

    getCallee(): Uint8Array | string;
    getCallee_asU8(): Uint8Array;
    getCallee_asB64(): string;
    setCallee(value: Uint8Array | string): void;

    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): void;

    getValue(): number;
    setValue(value: number): void;

    getGas(): number;
    setGas(value: number): void;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CallData.AsObject;
    static toObject(includeInstance: boolean, msg: CallData): CallData.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CallData, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CallData;
    static deserializeBinaryFromReader(message: CallData, reader: jspb.BinaryReader): CallData;
}

export namespace CallData {
    export type AsObject = {
        caller: Uint8Array | string,
        callee: Uint8Array | string,
        data: Uint8Array | string,
        value: number,
        gas: number,
    }
}
