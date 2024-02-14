import styled from "@emotion/styled";

export const InputWrap = styled.div`
  position: relative;
`;

export const Placeholder = styled.div`
  position: absolute;
  top: ${({ action }: { action: number }) => (action === 0 ? "0" : "-13px")};
  left: 0;

  text-align: left;

  font-size: ${({ action }: { action: number }) =>
    action === 0 ? ".9rem" : ".7rem"};

  transition: all 0.2s;
`;

export const Input = styled.input`
  border: none;
  border-bottom: 1px solid #555;

  width: 100%;

  outline: none;

  padding: 0.3rem;

  font-family: inherit;
`;
