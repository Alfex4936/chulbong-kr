import styled from "@emotion/styled";

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;

  position: absolute;

  top: 65px;
  right: 25px;

  background-color: #fff;

  border: 1px solid #ddd;
  border-radius: 9px;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;

  z-index: 10;
`;

export const InfoTop = styled.div`
  display: flex;
  align-items: center;

  border-bottom: 1px solid #eee;

  padding: 0.5rem;
`;

export const ProfileImgBox = styled.div`
  border-radius: 50%;

  box-shadow: rgba(0, 0, 0, 0.02) 0px 1px 3px 0px,
    rgba(27, 31, 35, 0.15) 0px 0px 0px 1px;

  width: 40px;
  height: 40px;

  margin-right: 0.5rem;

  overflow: hidden;

  & img {
    display: inline-block;

    width: 100%;
  }
`;

export const NameContainer = styled.div`
  margin-right: 1rem;

  text-align: left;

  font-weight: bold;

  & > div:last-of-type {
    width: 100%;

    font-size: 0.7rem;
    font-weight: 400;
    color: #555;
  }

  flex-grow: 1;
`;

export const LogoutButtonContainer = styled.div`
  & > button {
    font-size: 0.7rem;
  }
`;

export const InfoBottom = styled.div`
  display: flex;
`;
