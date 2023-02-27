import '@styles/globals.css';
import { SessionProvider, useSession } from 'next-auth/react';
import { StoreProvider } from '@components/providers/store';
import { useRouter } from 'next/router';
import { PayPalScriptProvider } from '@paypal/react-paypal-js';
import { Message } from '@components/ui/common';

const Noop = ({children}) => <>{children}</>

function MyApp({ Component, pageProps: { session, ...pageProps } }) {

    const Layout = Component.Layout ?? Noop

    return (
        <SessionProvider session={session}>
            <StoreProvider>
                <PayPalScriptProvider deferLoading={true}>
                    {Component.auth ? (
                        <Auth adminOnly={Component.auth.adminOnly}>
                            <Layout>
                                <Component {...pageProps} />
                            </Layout>
                        </Auth>
                    ) : (
                        <Layout>
                            <Component {...pageProps} />
                        </Layout>
                    )}
                </PayPalScriptProvider>
            </StoreProvider>
        </SessionProvider>
    );
}

function Auth({ children, adminOnly }) {
    const router = useRouter();
    const { status, data: session } = useSession({
        required: true,
        onUnauthenticated() {
            router.push('/unauthorized?message=login required');
        },
    });
    if (status === 'loading') {
        return <Message type="warning">
                            Is Loading...
                        </Message>
    }
    if (adminOnly && !session.user.isAdmin) {
        router.push('/unauthorized?message=admin login required');
    }

    return children;
}

export default MyApp;
