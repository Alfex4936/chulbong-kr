import styled from "@emotion/styled";

export const ImageBox = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  margin: auto;

  width: 200px;

  & button {
    background-color: #fff;
    border: none;
    cursor: pointer;

    &:hover > img {
      transform: scale(1.1);
    }
  }

  & img {
    display: inline-block;
    width: 100%;
  }
`;
