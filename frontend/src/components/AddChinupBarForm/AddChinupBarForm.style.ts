import styled from "@emotion/styled";

export const FormTitle = styled.h1`
  margin: 1rem;

  font-size: 1.5rem;
`;

export const InputWrap = styled.div`
  display: flex;
  flex-direction: column;

  margin-bottom: 1.7rem;

  & label {
    text-align: left;

    font-size: 0.8rem;
  }
`;

export const ButtonWrap = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;

  margin-top: 1rem;

  font-size: 0.8rem;

  color: #888;
`;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;
