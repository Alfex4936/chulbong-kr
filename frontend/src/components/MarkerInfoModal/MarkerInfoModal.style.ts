import styled from "@emotion/styled";

export const imageWrap = styled.div`
  & img {
    object-fit: cover;
    background-position: center;
    background-size: cover;

    border-radius: 1rem;
    box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
  }
`;

export const description = styled.div`
  font-size: 1.3rem;
  margin-top: 1rem;
`;
