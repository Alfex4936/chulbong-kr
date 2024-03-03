import styled from "@emotion/styled";

export const imageWrap = styled.div`
  position: relative;

  width: 90%;

  & img {
    object-fit: cover;
    background-position: center;
    background-size: cover;

    display: block;

    border-radius: 1rem;
    box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

    width: 100%;
  }

  margin: auto;
  margin-bottom: 2rem;
`;

export const description = styled.div`
  position: absolute;

  bottom: 0;
  left: 0;

  font-size: 1.3rem;
  margin-top: 1rem;
  color: #fff;

  white-space: nowrap;
  overflow: hidden;

  text-overflow: ellipsis;

  width: 100%;

  background-color: rgba(0, 0, 0, 0.5);

  padding: 1rem;

  border-bottom-left-radius: 1rem;
  border-bottom-right-radius: 1rem;

  &:hover {
    white-space: normal;
  }
`;

export const BottomButtons = styled.div`
  position: absolute;
  bottom: 1rem;
  left: 50%;
  transform: translateX(-50%);
`;
