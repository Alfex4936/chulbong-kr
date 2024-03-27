import styled from "@emotion/styled";

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  border: 1px solid red;

  width: 100%;
  height: 350px;

  margin-block: 1rem;

  border: 1.5px solid #888;
  border-radius: 0.5rem;
`;

export const JoinUser = styled.div`
  font-size: 0.8rem;

  color: #777;
`;

export const ConnectMessage = styled.div`
  font-size: 0.7rem;
  color: red;

  padding: 1rem;
`;

export const MessagesContainer = styled.div`
  padding: 1rem;

  overflow: auto;
`;

export const MessageWrapLeft = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;

  margin: 1rem;

  & > div:first-of-type {
    display: inline-block;

    padding: 0.3rem;
    background-color: #e9efff;
    border-radius: 4px;
  }
  & > div:last-of-type {
    font-size: 0.5rem;
    color: #666;
  }
`;

export const MessageWrapRight = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: flex-end;

  margin: 1rem;

  & > div:first-of-type {
    display: inline-block;

    padding: 0.3rem;
    background-color: #c0d0ff;
    border-radius: 4px;
  }
  & > div:last-of-type {
    font-size: 0.5rem;
    color: #666;
  }
`;

export const InputWrap = styled.div`
  display: flex;
  align-items: center;

  border-radius: 1rem;
  padding: 0 0.5rem;
  border: 1.5px solid #888;
  border-radius: 0.5rem;
`;

export const ReviewInput = styled.input`
  flex-grow: 1;

  border: none;
  outline: none;

  height: 1.5rem;

  font: inherit;
`;
