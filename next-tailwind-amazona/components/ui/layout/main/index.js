import { signOut, useSession } from "next-auth/react";
import Head from "next/head";
import Link from "next/link";
import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";
import { ToastContainer } from "react-toastify";
import { Menu } from "@headlessui/react";
import "react-toastify/dist/ReactToastify.css";
import { useStore } from "@components/providers";
import DropdownLink from "@components/ui/product/dropdown/DropdownLink";
import { useRouter } from "next/router";
import { SearchIcon } from "@heroicons/react/outline";
import { Footer } from "@components/ui/common";
import Image from "next/image";

export default function MainLayout({ title, children }) {
    const { status, data: session } = useSession();

    const { state, dispatch } = useStore();
    const { cart } = state;
    const [cartItemsCount, setCartItemsCount] = useState(0);

    useEffect(() => {
        setCartItemsCount(cart.cartItems.reduce((a, c) => a + c.quantity, 0));
    }, [cart.cartItems]);

    useEffect(() => {
        ///////////////////////////////////////
        // Slider
        const slider = () => {
            const slides = document.querySelectorAll(".custom-slide");
            const btnLeft = document.querySelector(".custom-slider__btn--left");
            const btnRight = document.querySelector(
                ".custom-slider__btn--right"
            );
            const dotContainer = document.querySelector(".dots");

            let curSlide = 0;
            const maxSlide = slides.length;

            // Functions
            const activateDot = function (slide) {
                document
                    .querySelectorAll(".dots__dot")
                    .forEach((dot) =>
                        dot.classList.remove("dots__dot--active")
                    );

                document
                    .querySelector(`.dots__dot[data-slide="${slide}"]`)
                    .classList.add("dots__dot--active");
            };

            const goToSlide = function (slide) {
                slides.forEach(
                    (s, i) =>
                        (s.style.transform = `translateX(${
                            100 * (i - slide)
                        }%)`)
                );
            };

            // Next slide
            const nextSlide = function () {
                if (curSlide === maxSlide - 1) {
                    curSlide = 0;
                } else {
                    curSlide++;
                }

                goToSlide(curSlide);
                activateDot(curSlide);
            };

            const prevSlide = function () {
                if (curSlide === 0) {
                    curSlide = maxSlide - 1;
                } else {
                    curSlide--;
                }
                goToSlide(curSlide);
                activateDot(curSlide);
            };

            const init = function () {
                goToSlide(0);
                activateDot(0);
            };
            init();

            // Event handlers
            btnRight.addEventListener("click", nextSlide);
            btnLeft.addEventListener("click", prevSlide);

            document.addEventListener("keydown", function (e) {
                if (e.key === "ArrowLeft") prevSlide();
                e.key === "ArrowRight" && nextSlide();
            });

            dotContainer.addEventListener("click", function (e) {
                if (e.target.classList.contains("dots__dot")) {
                    const { slide } = e.target.dataset;
                    goToSlide(slide);
                    activateDot(slide);
                }
            });
        };
        slider();
        console.log("runs")
    }, []);

    const logoutClickHandler = () => {
        Cookies.remove("cart");
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
                <title>{title ? title + " - Amazona" : "Amazona"}</title>
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
            <section id="hero" className="">
                {/* <!-- Hero Container --> */}
                <div className="container max-w-6xl mx-auto px-6 py-12">
                    {/* <!-- Menu/Logo Container --> */}
                    <nav className="flex items-center justify-between font-bold text-white">
                        {/* <!-- Logo --> */}
                        <div className="flex">
                            <Link href="/">
                                <a className="text-lg font-bold mr-5">
                                    Amazona
                                </a>
                            </Link>
                            <form
                                onSubmit={submitHandler}
                                className="mx-auto hidden w-full justify-center md:flex"
                            >
                                <input
                                    onChange={(e) => setQuery(e.target.value)}
                                    type="text"
                                    className="text-black rounded-tr-none rounded-br-none p-1 text-sm focus:ring-0"
                                    placeholder="Search products"
                                />
                                <button
                                    className="rounded rounded-tl-none rounded-bl-none bg-fuchsia-500 p-1 text-sm dark:text-black"
                                    type="submit"
                                    id="button-addon2"
                                >
                                    <SearchIcon className="h-5 w-5"></SearchIcon>
                                </button>
                            </form>
                        </div>
                        {/* <!-- Menu --> */}
                        <div className="hidden h-10 font-alata md:flex md:space-x-8">
                            <div className="group">
                                <Link href="/cart">
                                    <a className="p-2">
                                        Cart
                                        {cartItemsCount > 0 && (
                                            <span className="ml-1 rounded-full bg-red-600 px-2 py-1 text-xs font-bold text-white">
                                                {cartItemsCount}
                                            </span>
                                        )}
                                    </a>
                                </Link>
                                <div className="mx-2 group-hover:border-b group-hover:border-blue-50"></div>
                            </div>
                            <div className="group">
                                {status === "loading" ? (
                                    "Loading"
                                ) : session?.user ? (
                                    <Menu
                                        as="div"
                                        className="relative inline-block font-alata"
                                    >
                                        <Menu.Button className="font-alata font-bold">
                                            <div className="">
                                                Hi {session.user.name}!
                                            </div>
                                        </Menu.Button>
                                        <Menu.Items className="absolute right-0 w-56 origin-top-right border-2 shadow-lg ">
                                            <Menu.Item>
                                                <DropdownLink
                                                    className="dropdown-link"
                                                    href="/profile"
                                                >
                                                    Profile
                                                </DropdownLink>
                                            </Menu.Item>
                                            <Menu.Item>
                                                <DropdownLink
                                                    className="dropdown-link"
                                                    href="/order-history"
                                                >
                                                    Order History
                                                </DropdownLink>
                                            </Menu.Item>
                                            {session.user.isAdmin && (
                                                <Menu.Item>
                                                    <DropdownLink
                                                        className="dropdown-link"
                                                        href="/admin/dashboard"
                                                    >
                                                        Admin Dashboard
                                                    </DropdownLink>
                                                </Menu.Item>
                                            )}
                                            <Menu.Item>
                                                <a
                                                    className="dropdown-link"
                                                    href="#"
                                                    onClick={logoutClickHandler}
                                                >
                                                    Logout
                                                </a>
                                            </Menu.Item>
                                        </Menu.Items>
                                    </Menu>
                                ) : (
                                    <Link href="/login">
                                        <a className="p-2">Login</a>
                                    </Link>
                                )}
                                <div className="mx-2 group-hover:border-b group-hover:border-blue-50"></div>
                            </div>
                        </div>
                        {/* <!-- Hamburger Button --> */}
                        <div className="md:hidden">
                            <button
                                id="menu-btn"
                                type="button"
                                className="z-40 block hamburger md:hidden focus:outline-none"
                            >
                                <span className="hamburger-top"></span>
                                <span className="hamburger-middle"></span>
                                <span className="hamburger-bottom"></span>
                            </button>
                        </div>
                    </nav>

                    {/* <!-- Mobile Menu --> */}
                    <div
                        id="menu"
                        className="absolute top-0 bottom-0 left-0 hidden flex-col self-end w-full min-h-screen py-1 pt-40 pl-12 space-y-3 text-lg text-white uppercase bg-black"
                    >
                        <a href="#" className="hover:text-pink-500">
                            About
                        </a>
                        <a href="#" className="hover:text-pink-500">
                            Events
                        </a>
                        <a href="#" className="hover:text-pink-500">
                            Products
                        </a>
                        <a href="#" className="hover:text-pink-500">
                            Support
                        </a>
                    </div>

                    <div className="max-w-lg mt-32 mb-32 p-4 font-sans text-4xl text-white uppercase border-2 md:p-10 md:m-32 md:mx-0 md:text-6xl">
                        Impressive Technology That Deliver
                    </div>
                </div>
            </section>

            {/* <!-- Feature Section --> */}
            <section id="feature">
                {/* <!-- Feature Container --> */}
                <div className="relative container flex flex-col max-w-6xl mx-auto my-32 px-6 text-gray-900 md:flex-row md:px-0 font-sans">
                    {/* <!-- Image --> */}
                    <Image
                        src="/images/desktop/image-interactive.jpg"
                        alt=""
                        width="600"
                        height="500"
                    ></Image>

                    {/* <!-- Text Container --> */}
                    <div className="top-32 pr-0 bg-white md:absolute md:right-0 md:py-20 md:pl-20">
                        <h2 className="max-w-lg mt-10 mb-6 font-sans text-4xl text-center text-gray-900 uppercase md:text-5xl md:mt-0 md:text-left">
                            The leader in electronic commerce
                        </h2>

                        <p className="max-w-md text-center md:text-left">
                            Founded in 1994, Amazon has earned a reputation as a
                            disruptor of well-established industries through
                            technological innovation and aggressive reinvestment
                            of profits into capital expenditures. As of 2023, it
                            is the worlds largest online retailer and
                            marketplace.
                        </p>
                    </div>
                </div>
            </section>

            <div className="flex min-h-screen flex-col justify-between ">
                <div className="flex justify-center mb-20 md:justify-between">
                    <h2 className="ml-5 text-3xl text-center uppercase md:text-left md:text-4xl">
                        Our Products
                    </h2>
                </div>
                <main className="m-auto mt-4">{children}</main>
            </div>

            {/* Feedback */}
            <section className="section" id="section--3">
                <div className="section__title section__title--testimonials">
                    <h2 className="section__description">Not sure yet?</h2>
                    <h3 className="section__header">
                        Millions of Customers are already making their decisions.
                    </h3>
                </div>

                <div className="custom-slider">
                    <div className="custom-slide">
                        <div className="testimonial">
                            <h5 className="testimonial__header">
                                Best financial decision ever!
                            </h5>
                            <blockquote className="testimonial__text">
                            I can honestly say the Amazona is the single best investment I have ever made. Three years of using Amazona. But having loved it from the first day. I cannot imagine my life without it!
                            </blockquote>
                            <address className="testimonial__author">
                                <img
                                    src="img/user-1.jpg"
                                    alt=""
                                    className="testimonial__photo"
                                />
                                <h6 className="testimonial__name">
                                    Aarav Lynn
                                </h6>
                                <p className="testimonial__location">
                                    San Francisco, USA
                                </p>
                            </address>
                        </div>
                    </div>

                    <div className="custom-slide">
                        <div className="testimonial">
                            <h5 className="testimonial__header">
                                The last step to becoming a complete minimalist
                            </h5>
                            <blockquote className="testimonial__text">
                            Thank you @virmos and the @uet_vnu team -- 2 years with Amazona has been a LIFE changing experience - you have DEFINITELY succeeded in making ecommerce rocks - to the future #ai #machinelearning #amazona
                            </blockquote>
                            <address className="testimonial__author">
                                <img
                                    src="img/user-2.jpg"
                                    alt=""
                                    className="testimonial__photo"
                                />
                                <h6 className="testimonial__name">
                                    Miyah Miles
                                </h6>
                                <p className="testimonial__location">
                                    London, UK
                                </p>
                            </address>
                        </div>
                    </div>

                    <div className="custom-slide">
                        <div className="testimonial">
                            <h5 className="testimonial__header">
                                Finally free from old-school ecommerce website
                            </h5>
                            <blockquote className="testimonial__text">
                            O meu pai chama-se Miguel. Tem 58 anos, é alto e moreno, com olhos castanhos. Conheceu a minha mãe, Maria, na faculdade, quando estavam a estudar psicologia. A minha mãe tem 55 anos e é ruiva, com olhos azuis. A minha irmã chama-se Joana e tem 25 anos. Tem o cabelo ruivo, como a minha mãe, e os olhos castanhos como o meu pai. 
                            </blockquote>
                            <address className="testimonial__author">
                                <img
                                    src="img/user-3.jpg"
                                    alt=""
                                    className="testimonial__photo"
                                />
                                <h6 className="testimonial__name">
                                    Francisco Gomes
                                </h6>
                                <p className="testimonial__location">
                                    Lisbon, Portugal
                                </p>
                            </address>
                        </div>
                    </div>

                    <button className="custom-slider__btn custom-slider__btn--left">
                        &larr;
                    </button>
                    <button className="custom-slider__btn custom-slider__btn--right">
                        &rarr;
                    </button>
                    <div className="dots">
                        <button className="dots__dot" data-slide="0"></button>
                        <button className="dots__dot" data-slide="1"></button>
                        <button className="dots__dot" data-slide="2"></button>
                    </div>
                </div>
            </section>
            <Footer />
        </>
    );
}
