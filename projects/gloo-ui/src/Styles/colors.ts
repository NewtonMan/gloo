const colorsList = {
  kumquatOrange: 'hsl(22, 100%, 94%)', //#FFEADE
  tangerineOrange: 'hsl(20, 92%, 84%)', //#FCCCB3
  peachOrange: 'hsl(20, 75%, 64%)', //#E98D5E
  grapefruitOrange: 'hsl(20, 76%, 47%)', //#D75B1D
  pumpkinOrange: 'hsl(19, 100%, 27%)', //#8B2E00

  flashlightGold: 'hsl(39, 100%, 95%)', //#FFF6E5
  sunGold: 'hsl(39, 100%, 59%)', //#FFB831,
  lightGold: 'hsl(39, 100%, 84%)',
  darkGold: 'hsl(39, 100%, 20%)',

  groveGreen: 'hsl(161, 100%, 96%)', //	#ebfff9
  forestGreen: 'hsl(161, 88%, 42%)', //#0dc98e

  dropBlue: 'hsl(202, 100%, 98%)', //#F7FCFF
  splashBlue: 'hsla(198, 100%, 93%, 1)', //#DEF5FF
  puddleBlue: 'hsl(198, 100%, 84%)', //#B1E7FF
  pondBlue: 'hsl(198, 81%, 67%)', // #6AC7F0
  lakeBlue: 'hsl(198, 71%, 60%)', //#54B7E3
  seaBlue: 'hsl(198, 72%, 46%)', //#2196C9
  oceanBlue: 'hsl(198, 84%, 37%)', //##0F7FB1
  blue100: '#F7FCFF' as '#F7FCFF',
  blue200: '#EAF8FF' as '#EAF8FF',
  blue300: '#B1E7FF' as '#B1E7FF',
  blue400: '#6AC7F0' as '#6AC7F0',
  blue500: '#54B7E3' as '#54B7E3',
  blue600: '#2196C9' as '#2196C9',
  blue700: '#0F7FB1' as '#0F7FB1',
  blue800: '#10344C' as '#10344C',
  blue900: '#0e2636' as '#0e2636',

  januaryGrey: 'hsl(0, 0%, 97%)', //#F9F9F9',
  februaryGrey: 'hsl(0, 0%, 94%)', //#F2F2F2',
  marchGrey: 'hsl(240, 1%, 87%)', //#DFDFE0
  aprilGrey: 'hsl(216, 13%, 85%)', //#D4D8DE
  mayGrey: 'hsl(214, 12%, 78%)', //#C2C8D0
  juneGrey: 'hsl(216, 10%, 70%)', //#ADB3BC
  septemberGrey: 'hsl(200, 3%, 44%)', //#6E7477
  novemberGrey: 'hsl(200, 5%, 22%)', //#35393B

  scrollbarBorderGrey: '#E5E5E5',
  scrollbarBackgroundGrey: 'hsl(0,0%,98.5%)',

  boxShadow: 'hsla(0, 0%, 0%, 0.1)',
  darkerBoxShadow: 'hsla(0,0%,0%, .15)'
};
type ColorProps = typeof colorsList;

// this gets us autocomplete!
export const colors = { ...colorsList } as ColorProps;

export function hslToHSLA(hsl: string, alpha: number | string) {
  let hslContent = hsl.split(')')[0].split('(')[1];

  return `hsla(${hslContent}, ${alpha})`;
}
