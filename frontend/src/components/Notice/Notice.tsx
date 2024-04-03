import * as Styled from "./Notice.style";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import NotificationsNoneIcon from "@mui/icons-material/NotificationsNone";

const Notice = () => {
  return (
    <Styled.Container>
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel1-content"
          id="panel1-header"
        >
          <span>
            <NotificationsNoneIcon />
          </span>
          <Styled.NoticeTitle>위치별 채팅 기능 추가</Styled.NoticeTitle>
        </AccordionSummary>
        <AccordionDetails>
          <Styled.NoticeContent>
            철봉에 위치별로 익명 채팅이 가능합니다.
          </Styled.NoticeContent>
          <Styled.ImgWrap>
            <img src="/notice/1.png" />
          </Styled.ImgWrap>
        </AccordionDetails>
      </Accordion>
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel2-content"
          id="panel2-header"
        >
          <span>
            <NotificationsNoneIcon />
          </span>
          <Styled.NoticeTitle>지역별 채팅 기능 추가</Styled.NoticeTitle>
        </AccordionSummary>
        <AccordionDetails>
          <Styled.NoticeContent>
            현재 지도의 위치에 따라 지역별 채팅이 가능합니다.
          </Styled.NoticeContent>
          <Styled.ImgWrap>
            <img src="/notice/2.png" />
          </Styled.ImgWrap>
        </AccordionDetails>
      </Accordion>
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel2-content"
          id="panel2-header"
        >
          <span>
            <NotificationsNoneIcon />
          </span>
          <Styled.NoticeTitle>철봉 위치 날씨 표시</Styled.NoticeTitle>
        </AccordionSummary>
        <AccordionDetails>
          <Styled.NoticeContent>
            철봉 위치 정보에 날씨 정보가 포함됩니다.
          </Styled.NoticeContent>
          <Styled.ImgWrap>
            <img src="/notice/3.png" />
          </Styled.ImgWrap>
        </AccordionDetails>
      </Accordion>
    </Styled.Container>
  );
};

export default Notice;
