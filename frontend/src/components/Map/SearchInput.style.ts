import styled from "@emotion/styled";

export const InputWrap = styled.div`
  display: flex;
  align-items: center;

  padding: 0 0.5rem;
  margin-right: 0.5rem;

  border-radius: 0.5rem;

  background-color: #fff;

  z-index: 200;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
`;

export const SearchInput = styled.input`
  flex-grow: 1;

  border: none;
  outline: none;

  width: 100%;

  height: 1.5rem;

  font: inherit;

  background-color: transparent;
`;

export const Result = styled.div`
  position: absolute;

  top: 45px;
  left: 0;

  border-radius: 0.5rem;

  background-color: #fff;

  width: 100%;

  padding: 1rem 0;
  margin-right: 0.5rem;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
`;

export const ResultItem = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;
  margin: 0 auto 1rem auto;

  width: 90%;

  border-radius: 0.4rem;
  background-color: #e9efff;

  & > div:first-of-type {
    flex-grow: 1;

    & > span:last-of-type {
      font-size: 0.7rem;
      color: #777;
    }
  }
`;
