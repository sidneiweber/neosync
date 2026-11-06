import { cn } from '@/libs/utils';
import { ReactElement } from 'react';

interface MaskSyncSymbolProps {
  className?: string;
}

// The MaskSync symbol: a solid triangle (the original data) and its striped twin
// (the masked copy). Uses currentColor so callers set the color with text-* classes.
export function MaskSyncSymbol(props: MaskSyncSymbolProps): ReactElement {
  const { className } = props;
  return (
    <svg
      viewBox="14 14 36 36"
      fill="currentColor"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
      className={className}
    >
      <polygon points="14,14 48,14 14,48" />
      <polygon points="48,18 50,18 50,22 44,22" />
      <polygon points="41,25 50,25 50,29 37,29" />
      <polygon points="34,32 50,32 50,36 30,36" />
      <polygon points="27,39 50,39 50,43 23,43" />
      <polygon points="20,46 50,46 50,50 16,50" />
    </svg>
  );
}

interface Props {
  className?: string;
}

// Brand signature: symbol + lowercase wordmark. Ink symbol on light backgrounds,
// Mask Lime symbol on dark backgrounds.
export function MainLogo(props: Props): ReactElement {
  const { className } = props;
  return (
    <span
      className={cn('flex items-center gap-2', className)}
      aria-label="MaskSync"
    >
      <MaskSyncSymbol className="h-6 w-6 text-gray-950 dark:text-lime" />
      <span
        className="font-display text-xl font-semibold leading-none text-foreground"
        style={{ letterSpacing: '-0.03em' }}
      >
        masksync
      </span>
    </span>
  );
}
