import { DropdownLink } from "@components/ui/product";
import { Menu } from "@headlessui/react";
import { SearchIcon } from "@heroicons/react/outline";
import Link from "next/link";

export default function HeroSection({status, session, cartItemsCount, submitHandler, logoutClickHandler, setQuery}) {
    return (
        <section id="hero" className="">
            {/* <!-- Hero Container --> */}
            <div className="container max-w-6xl mx-auto px-6 py-12">
                {/* <!-- Menu/Logo Container --> */}
                <nav className="flex items-center justify-between font-bold text-white">
                    {/* <!-- Logo --> */}
                    <div className="flex">
                        <Link href="/">
                            <a className="text-lg font-bold mr-5">
                                Shopee
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
    );
}
