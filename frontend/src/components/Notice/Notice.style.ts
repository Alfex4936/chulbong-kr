import styled from "@emotion/styled";

export const Container = styled.div`
  max-height: 400px;

  overflow: auto;
`;

export const TitleIcon = styled.span`
  margin-right: 1.5rem;
`;

export const NoticeTitle = styled.h1`
  font-weight: bold;

  font-size: 1rem;
`;

export const NoticeContent = styled.div`
  font-weight: bold;

  margin-bottom: 1.3rem;
`;

export const ImgWrap = styled.div`
  & > img {
    box-shadow: rgba(99, 99, 99, 0.2) 0px 2px 8px 0px;

    max-width: 85%;
  }
  & > img:hover {
    box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
  }
`;
