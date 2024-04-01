import ReplyIcon from "@mui/icons-material/Reply";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { Fragment, useEffect, useRef, useState } from "react";
import useAddressData from "../../hooks/useAddressData";
import useInput from "../../hooks/useInput";
import useChatIdStore from "../../store/useChatIdStore";
import getRegion from "../../utils/getRegionCode";
import type { ChatMessage, Chatdata } from "../MarkerInfoModal/ChatRoom";
import * as Styled from "./LocalChat.style";

interface Props {
  setLocalChat: React.Dispatch<React.SetStateAction<boolean>>;
}

const LocalChat = ({ setLocalChat }: Props) => {
  const cidState = useChatIdStore();

  const chatValue = useInput("");

  const { address, isError } = useAddressData();

  const ws = useRef<WebSocket | null>(null);
  const chatBox = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [connection, setConnection] = useState(false);

  const [messages, setMessages] = useState<Chatdata[]>([]);
  const [connectionMsg, setConnectionMsg] = useState("");

  const [roomTitle, setRoomTitle] = useState("");

  useEffect(() => {
    const code = getRegion(address?.depth1 as string).getCode();
    if (!address || isError || code === "") {
      setConnectionMsg("채팅 서비스를 지원하지 않는 지역입니다!");
      return;
    }

    ws.current = new WebSocket(
      `wss://api.k-pullup.com/ws/${code}?request-id=${cidState.cid}`
    );

    ws.current.onopen = () => {
      setConnection(true);
      setConnectionMsg(
        "비속어 사용에 주의해주세요. 이후 서비스 사용이 제한될 수 있습니다!"
      );
    };

    ws.current.onmessage = async (event) => {
      const data: ChatMessage = JSON.parse(event.data);
      if (data.userNickname === "chulbong-kr") {
        const titleArr = data.message.split(" ");

        titleArr[0] = getRegion(data.roomID).getTitle();

        setRoomTitle(titleArr.join(" "));
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
      setLocalChat(false);
    };

    ws.current.onclose = () => {
      setLocalChat(false);
      console.log("연결 종료");
    };

    return () => {
      ws.current?.close();
    };
  }, [address]);

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
    chatValue.reset();
    inputRef.current?.focus();
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      handleChat();
    }
  };

  return (
    <div>
      <div>{roomTitle}</div>
      <Styled.Container>
        <div>
          <Styled.ConnectMessage>
            {connection ? connectionMsg : "채팅방 접속중..."}
          </Styled.ConnectMessage>
        </div>
        <Styled.MessagesContainer ref={chatBox}>
          {messages.map((message) => {
            if (message.name === "chulbong-kr") return;
            if (message.msg?.includes("님이 입장하셨습니다.")) {
              return (
                <Styled.JoinUser key={message.mid}>
                  {message.name}님이 참여하였습니다.
                </Styled.JoinUser>
              );
            }
            if (message.msg?.includes("님이 퇴장하셨습니다.")) {
              return (
                <Styled.JoinUser key={message.mid}>
                  {message.name}님이 나가셨습니다.
                </Styled.JoinUser>
              );
            }
            return (
              <Fragment key={message.mid}>
                {message.userid === cidState.cid ? (
                  <Styled.MessageWrapRight>
                    <div>{message.msg}</div>
                    <div>{message.name}</div>
                  </Styled.MessageWrapRight>
                ) : (
                  <Styled.MessageWrapLeft>
                    <div>{message.msg}</div>
                    <div>{message.name}</div>
                  </Styled.MessageWrapLeft>
                )}
              </Fragment>
            );
          })}
        </Styled.MessagesContainer>
      </Styled.Container>
      <Styled.InputWrap>
        <Styled.ReviewInput
          maxLength={40}
          ref={inputRef}
          disabled={!connection}
          type="text"
          name="reveiw-content"
          value={chatValue.value}
          onChange={chatValue.onChange}
          onKeyDown={handleKeyPress}
        />
        <Tooltip title="보내기" arrow disableInteractive>
          <IconButton onClick={handleChat} aria-label="send">
            <ReplyIcon />
          </IconButton>
        </Tooltip>
      </Styled.InputWrap>
    </div>
  );
};

export default LocalChat;
