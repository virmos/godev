import { signOut, useSession } from "next-auth/react";
import Head from "next/head";
import React, { useEffect, useState } from "react";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { useStore } from "@components/providers";
import { useRouter } from "next/router";
import { Footer } from "@components/ui/common";
import HeroSection from "../section/hero";
import FeatureSection from "../section/feature";
import FeedbackSection from "../section/feedback";
import slider from "@utils/animation";

export default function MainLayout({ title, children }) {
    const { status, data: session } = useSession();

    const { state, dispatch } = useStore();
    const { cart } = state;
    const [cartItemsCount, setCartItemsCount] = useState(0);

    useEffect(() => {
        setCartItemsCount(cart.cartItems.reduce((a, c) => a + c.quantity, 0));
    }, [cart.cartItems]);

    useEffect(() => {
        // prepare sliding animation
        slider();
    }, []);

    const logoutClickHandler = () => {
        dispatch({ type: "CART_RESET" });
        signOut({ callbackUrl: "/login" });
    };

    const [query, setQuery] = useState("");

    const router = useRouter();
    const submitHandler = (e) => {
        e.preventDefault();
        router.push(`/search?query=${query}`);
    };

    return (
        <>
            <Head>
                <title>{title ? title + " - Shopee" : "Shopee"}</title>
                <meta name="description" content="Ecommerce Website" />
                <link rel="icon" href="/favicon.ico" />
                <link rel="preconnect" href="https://fonts.googleapis.com" />
                <link
                    rel="preconnect"
                    href="https://fonts.gstatic.com"
                    crossOrigin
                />
                <link
                    href="https://fonts.googleapis.com/css2?family=Alata&family=Josefin+Sans:wght@300&display=swap"
                    rel="stylesheet"
                />
            </Head>

            <ToastContainer position="bottom-center" limit={1} />

            {/* <!-- Hero Section --> */}
            <HeroSection
                status={status}
                session={session}
                cartItemsCount={cartItemsCount}
                submitHandler={submitHandler}
                logoutClickHandler={logoutClickHandler}
                setQuery={setQuery}
            />

            {/* <!-- Feature Section --> */}
            <FeatureSection/>

            <div className="flex min-h-screen flex-col justify-between ">
                <div className="flex justify-center mb-20 md:justify-between">
                    <h2 className="ml-5 text-3xl text-center uppercase md:text-left md:text-4xl">
                        Our Products
                    </h2>
                </div>
                <main className="m-auto mt-4">{children}</main>
            </div>

            {/* Feedback */}
            <FeedbackSection/>
            <Footer />
        </>
    );
}
