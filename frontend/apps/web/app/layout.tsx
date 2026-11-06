import '@/app/globals.css';
import { GoogleScriptProvider } from '@/components/providers/googleTag-provider';
import {
  PHProvider,
  PostHogPageview,
} from '@/components/providers/posthog-provider';
import { ThemeProvider } from '@/components/providers/theme-provider';
import { UnifyScriptProvider } from '@/components/providers/unify-provider';
import { Metadata } from 'next';
import { Inter, JetBrains_Mono, Space_Grotesk } from 'next/font/google';
import { ReactElement, Suspense } from 'react';
import BaseLayout from './BaseLayout';

const inter = Inter({
  subsets: ['latin', 'latin-ext'],
  variable: '--font-inter',
  display: 'swap',
});
const spaceGrotesk = Space_Grotesk({
  subsets: ['latin', 'latin-ext'],
  weight: ['500', '600'],
  variable: '--font-space-grotesk',
  display: 'swap',
});
const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin', 'latin-ext'],
  variable: '--font-jetbrains-mono',
  display: 'swap',
});

export const metadata: Metadata = {
  title: 'MaskSync',
  description:
    'Open source data anonymization, synthetic data and environment sync.',
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}): Promise<ReactElement> {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${inter.variable} ${spaceGrotesk.variable} ${jetbrainsMono.variable}`}
    >
      <head />
      <body className="min-h-screen bg-background font-sans antialiased overflow-scroll">
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <>
            <PHProvider>
              <BaseLayout>
                <>
                  <Suspense>
                    <UnifyScriptProvider />
                    <GoogleScriptProvider />
                  </Suspense>
                  <Suspense>
                    <PostHogPageview />
                  </Suspense>
                  {children}
                </>
              </BaseLayout>
            </PHProvider>
          </>
        </ThemeProvider>
      </body>
    </html>
  );
}
