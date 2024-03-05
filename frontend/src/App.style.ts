import styled from "@emotion/styled";

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;

export const ChangePasswordButtonsWrap = styled.div`
  display: flex;

  & button {
    margin: 1rem 0.5rem 0 0.5rem;
  }
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
