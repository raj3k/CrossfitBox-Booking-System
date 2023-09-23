import {Button, Input} from '@mui/joy';
import * as React from 'react';
import { styled } from '@mui/joy/styles';
import { useParams } from 'react-router-dom';
import { useState } from 'react';
import toast from 'react-hot-toast';
import * as api from "../helpers/api";

const StyledInput = styled('input')({
  border: 'none', // remove the native input border
  minWidth: 0, // remove the native input width
  outline: 0, // remove the native input outline
  padding: 0, // remove the native input padding
  paddingTop: '1em',
  flex: 1,
  color: 'inherit',
  backgroundColor: 'transparent',
  fontFamily: 'inherit',
  fontSize: 'inherit',
  fontStyle: 'inherit',
  fontWeight: 'inherit',
  lineHeight: 'inherit',
  textOverflow: 'ellipsis',
  '&::placeholder': {
    opacity: 0,
    transition: '0.1s ease-out',
  },
  '&:focus::placeholder': {
    opacity: 1,
  },
  '&:focus ~ label, &:not(:placeholder-shown) ~ label, &:-webkit-autofill ~ label': {
    top: '0.5rem',
    fontSize: '0.75rem',
  },
  '&:focus ~ label': {
    color: 'var(--Input-focusedHighlight)',
  },
  '&:-webkit-autofill': {
    alignSelf: 'stretch', // to fill the height of the root slot
  },
  '&:-webkit-autofill:not(* + &)': {
    marginInlineStart: 'calc(-1 * var(--Input-paddingInline))',
    paddingInlineStart: 'var(--Input-paddingInline)',
    borderTopLeftRadius:
      'calc(var(--Input-radius) - var(--variant-borderWidth, 0px))',
    borderBottomLeftRadius:
      'calc(var(--Input-radius) - var(--variant-borderWidth, 0px))',
  },
});

const StyledLabel = styled('label')(({ theme }) => ({
  position: 'absolute',
  lineHeight: 1,
  top: 'calc((var(--Input-minHeight) - 1em) / 2)',
  color: theme.vars.palette.text.tertiary,
  fontWeight: theme.vars.fontWeight.md,
  transition: 'all 150ms cubic-bezier(0.4, 0, 0.2, 1)',
}));

const InnerInput = React.forwardRef<
  HTMLInputElement,
  JSX.IntrinsicElements['input']
>(function InnerInput(props, ref) {
  const id = React.useId();
  return (
    <React.Fragment>
      <StyledInput {...props} ref={ref} id={id} />
      <StyledLabel htmlFor={id}>Token</StyledLabel>
    </React.Fragment>
  );
});

const Activate: React.FC = () => {
    const { userId } = useParams();
    const [token, setToken] = useState("");
    const allowActivate = token.length > 0;

    const handleTokenInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
      const text = e.target.value as string;
      setToken(text);
    };

    const handleActivateBtnClick = async (e: React.FormEvent) => {
      e.preventDefault();
      if (userId && token) {
        try {
          let tokenVal = token.split(' ').join('');
          const finalVal = token.match(/.{1,3}/g)!.join(' ');
          tokenVal = finalVal;
          
          await api.activateUser(userId, tokenVal)
        } catch(error: any) {
          console.error(error);
          toast.error(error.response.data.message);
        }
      } else {
        toast.error("userId is undefined")
      }
    }

    return (
      <div className="flex flex-row justify-center items-center w-full h-auto mt-12 sm:mt-24 bg-white">
        <div className="w-80 max-w-full h-full py-4 flex flex-col justify-start items-center">
          <div className="w-full py-4 grow flex flex-col justify-center items-center">
            <p className="w-full text-2xl mt-6 flex flex-col justify-center items-center">Activate your account</p>
            <form className="w-full mt-4" onSubmit={handleActivateBtnClick}>
              <div className="flex flex-col justify-start items-start w-full gap-4 py-4">
              <Input
                className="w-full"
                size="lg"
                type="text"
                placeholder="Token"
                onChange={handleTokenInputChanged}
                required
                slots={{ input: InnerInput }}
                slotProps={{ input: { placeholder: 'XXX XXX', type: 'text' } }}
                sx={{
                  '--Input-minHeight': '56px',
                  '--Input-radius': '6px',
                }}
              />
              </div>
              <div className="flex flex-col justify-center items-center w-full mt-6">
                <Button
                  type="submit"
                  color="primary"
                  disabled={!allowActivate}
                >
                  Activate
                </Button>
            </div>
            </form>
          </div>
        </div>
      </div>
    )
}

export default Activate;