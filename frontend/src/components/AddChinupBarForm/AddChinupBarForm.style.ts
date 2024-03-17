import styled from "@emotion/styled";

export const FormTitle = styled.h1`
  margin: 1rem;

  font-size: 1.5rem;
`;

export const FlexCenter = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Empty = styled.div`
  flex-grow: 1;
`;

export const Count = styled.div`
  margin: 0 .5rem;
  width: 17px;
`;

export const NumberInputWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;

  border: 0.5px solid #888;
  border-radius: 0.3rem;
  padding: 1rem;

  margin-bottom: 1rem;

  & input {
    border: 1px solid red;
    width: 40px;
  }

  & > div {
    display: flex;
    width: 100%;
  }
  // & label {
  //   text-align: left;

  //   font-size: 0.8rem;
  // }

  // & > div:first-of-type {
  //   margin-right: 3rem;
  // }
  // & > div:last-of-type {
  //   margin-left: 3rem;
  // }
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
