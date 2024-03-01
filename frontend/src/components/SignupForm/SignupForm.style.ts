import styled from "@emotion/styled";

export const FormTitle = styled.h1`
  margin: 1rem;

  font-size: 1.5rem;
`;

export const InputWrap = styled.div`
  position: relative;

  display: flex;
  flex-direction: column;

  margin-bottom: 1.7rem;

  & label {
    text-align: left;

    font-size: 0.8rem;
  }
`;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;

export const TimerContainer = styled.div`
  font-size: 0.8rem;
`;
