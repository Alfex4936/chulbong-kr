import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import ReplyIcon from "@mui/icons-material/Reply";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { Fragment, useEffect, useRef, useState } from "react";
import useInput from "../../hooks/useInput";
import useChatIdStore from "../../store/useChatIdStore";
import * as Styled from "./ChatRoom.style";

interface ChatMessage {
  uid: string;
  message: string;
  userId: string;
  userNickname: string;
  roomID: string;
  timestamp: number;
  isOwner: boolean;
}

interface Chatdata {
  msg: string;
  name: string;
  isOwner: boolean;
  mid: string;
  userid: string;
}

interface Props {
  setIsChatView: React.Dispatch<React.SetStateAction<boolean>>;
  markerId: number;
}

const ChatRoom = ({ setIsChatView, markerId }: Props) => {
  const cidState = useChatIdStore();

  const chatValue = useInput("");

  const ws = useRef<WebSocket | null>(null);
  const chatBox = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [connection, setConnection] = useState(false);

  const [messages, setMessages] = useState<Chatdata[]>([]);
  const [connectionMsg, setConnectionMsg] = useState("");

  useEffect(() => {
    ws.current = new WebSocket(
      `wss://api.k-pullup.com/ws/${markerId}?request-id=${cidState.cid}`
    );

    ws.current.onopen = () => {
      setConnection(true);
      setConnectionMsg(
        "비속어 사용에 주의해주세요. 이후 서비스 사용이 제한될 수 있습니다!"
      );
    };

    ws.current.onmessage = async (event) => {
      const data: ChatMessage = JSON.parse(event.data);

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
    };

    ws.current.onclose = () => {
      setIsChatView(false);
      console.log("연결 종료");
    };

    return () => {
      ws.current?.close();
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
      <Tooltip title="이전" arrow disableInteractive>
        <IconButton
          onClick={() => {
            setIsChatView(false);
          }}
          aria-label="delete"
          sx={{
            position: "absolute",
            top: "0",
            left: "0",
          }}
        >
          <ArrowBackIcon />
        </IconButton>
      </Tooltip>
      <Styled.Container>
        <div>
          <Styled.ConnectMessage>
            {connection ? connectionMsg : "채팅방 접속중..."}
          </Styled.ConnectMessage>
        </div>
        <Styled.MessagesContainer ref={chatBox}>
          {messages.map((message) => {
            if (message.msg.includes("님이 입장하셨습니다.")) {
              return (
                <Styled.JoinUser key={message.mid}>
                  {message.name}님이 참여하였습니다.
                </Styled.JoinUser>
              );
            }
            if (message.msg.includes("님이 퇴장하셨습니다.")) {
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

export default ChatRoom;
