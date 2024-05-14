"use client";

import Heading from "@/components/atom/Heading";
import SendIcon from "@/components/icons/SendIcon";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useChatIdStore from "@/store/useChatIdStore";
import { Fragment, useEffect, useRef, useState } from "react";

export interface ChatMessage {
  uid: string;
  message: string;
  userId: string;
  userNickname: string;
  roomID: string;
  timestamp: number;
  isOwner: boolean;
}

export interface Chatdata {
  msg: string;
  name: string;
  isOwner: boolean;
  mid: string;
  userid: string;
}

interface Props {
  markerId: number;
}

const PullupChatClient = ({ markerId }: Props) => {
  const cidState = useChatIdStore();

  const chatValue = useInput("");

  const ws = useRef<WebSocket | null>(null);
  const chatBox = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [connection, setConnection] = useState(false);

  const [messages, setMessages] = useState<Chatdata[]>([]);
  const [connectionMsg, setConnectionMsg] = useState("");

  const [roomTitle, setRoomTitle] = useState("");
  const [roomSubTitle, setRoomSubTitle] = useState("");

  const [isChatError, setIsChatError] = useState(false);

  useEffect(() => {
    ws.current?.close();

    ws.current = new WebSocket(
      `wss://api.k-pullup.com/ws/${markerId}?request-id=${cidState.cid}`
    );

    ws.current.onopen = () => {
      setMessages([]);
      setConnection(true);
      setConnectionMsg(
        "비속어 사용에 주의해주세요. 이후 서비스 사용이 제한될 수 있습니다!"
      );
    };

    ws.current.onmessage = async (event) => {
      const data: ChatMessage = JSON.parse(event.data);
      if (data.userNickname === "chulbong-kr") {
        setRoomTitle(data.message);
        console.log(data.message);
      }

      setMessages((prevMessages) => [
        ...prevMessages,
        {
          msg: data.message,
          name: data.userNickname,
          isOwner: data.isOwner,
          mid: data.uid,
          userid: data.userId,
        },
      ]);
    };

    ws.current.onerror = (error) => {
      setConnectionMsg(
        "채팅방에 참여 중 에러가 발생하였습니다. 잠시 후 다시 시도해 주세요!"
      );
      console.error("연결 에러:", error);
      setIsChatError(true);
    };

    ws.current.onclose = () => {
      console.log("연결 종료");
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!ws) return;
    const pingInterval = setInterval(() => {
      ws.current?.send(JSON.stringify({ type: "ping" }));
    }, 30000);

    return () => {
      clearInterval(pingInterval);
    };
  }, []);

  useEffect(() => {
    const scrollBox = chatBox.current;

    if (scrollBox) {
      scrollBox.scrollTop = scrollBox.scrollHeight;
    }
  }, [messages]);

  const handleChat = () => {
    if (chatValue.value === "") return;
    ws.current?.send(chatValue.value);
    chatValue.resetValue();
    inputRef.current?.focus();
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      handleChat();
    }
  };

  if (isChatError) {
    return (
      <div>
        <Heading title={`채팅방`} />
        <div className="text-red text-center">
          채팅을 불러오는데 실패하였습니다. <br /> 잠시 후 다시 시도해 주세요.
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      <Heading title={roomTitle} subTitle={roomSubTitle} />
      <div
        className="grow w-full flex flex-col justify-between px-3"
        ref={chatBox}
      >
        <div className="text-center text-red text-sm">
          {connection ? connectionMsg : "채팅방 접속중..."}
        </div>
        <div>
          {messages.map((message) => {
            if (message.name === "chulbong-kr") return;
            if (message.msg?.includes("님이 입장하셨습니다.")) {
              return (
                <div
                  key={message.mid}
                  className="truncate px-5 py-2 text-center text-sm text-grey-dark"
                >
                  {message.name}님이 참여하였습니다.
                </div>
              );
            }
            if (message.msg?.includes("님이 퇴장하셨습니다.")) {
              return (
                <div
                  key={message.mid}
                  className="truncate px-5 py-2 text-center text-sm text-grey-dark"
                >
                  {message.name}님이 나가셨습니다.
                </div>
              );
            }
            return (
              <Fragment key={message.mid}>
                {message.userid === cidState.cid ? (
                  <div className="flex flex-col items-end w-full">
                    <div className="max-w-1/2 p-1 rounded-lg bg-slate-700 shadow-sm">
                      {message.msg}
                    </div>
                    <div className="text-xs text-grey-dark">{message.name}</div>
                  </div>
                ) : (
                  <div className="flex flex-col items-start w-full">
                    <div className="max-w-1/2 p-1 rounded-lg bg-slate-600 shadow-sm">
                      {message.msg}
                    </div>
                    <div className="text-xs text-grey-dark">{message.name}</div>
                  </div>
                )}
              </Fragment>
            );
          })}
        </div>
      </div>
      <div className="relative flex items-center justify-center w-full h-14 px-3">
        <Input
          type="text"
          ref={inputRef}
          maxLength={40}
          disabled={!connection}
          name="reveiw-content"
          value={chatValue.value}
          onChange={chatValue.handleChange}
          onKeyDown={handleKeyPress}
          className="bg-black-light-2 text-base pr-8"
        />
        <button className="absolute right-5" onClick={handleChat}>
          <SendIcon size={20} />
        </button>
      </div>
    </div>
  );
};

export default PullupChatClient;
