import Image from 'next/image';

export default function Footer() {

    return (
        <footer className="bg-black font-sans">
            {/* <!-- Container --> */}
            <div className="container max-w-6xl py-10 mx-auto">
                <div
                    className="flex flex-col items-center mb-8 space-y-6 md:flex-row md:space-y-0 md:justify-between md:items-start"
                >
                    {/* <!-- Menu & Logo Container --> */}
                    <div
                        className="flex flex-col items-center space-y-8 md:items-start md:space-y-4"
                    >
                        {/* <!-- Logo -->
                        <div className="">
                            <Image
                                src="/logo.png"
                                alt="icon-logo"
                                width={100}
                                height={100}
                            ></Image>
                        </div> */}
                        {/* <!-- Menu Container --> */}
                        <div
                            className="flex flex-col items-center space-y-4 font-bold text-white md:flex-row md:space-y-0 md:space-x-6 md:ml-3"
                        >
                            {/* <!-- Item 1 --> */}
                            <div className="h-10 group">
                                <a href="#">About</a>
                                <div
                                    className="mx-2 group-hover:border-b group-hover:border-blue-25"
                                ></div>
                            </div>
                            {/* <!-- Item 2 --> */}
                            <div className="h-10 group">
                                <a href="#">Events</a>
                                <div
                                    className="mx-2 group-hover:border-b group-hover:border-blue-25"
                                ></div>
                            </div>
                            {/* <!-- Item 3 --> */}
                            <div className="h-10 group">
                                <a href="#">Products</a>
                                <div
                                    className="mx-2 group-hover:border-b group-hover:border-blue-25"
                                ></div>
                            </div>
                            {/* <!-- Item 4 --> */}
                            <div className="h-10 group">
                                <a href="#">Support</a>
                                <div
                                    className="mx-2 group-hover:border-b group-hover:border-blue-25"
                                ></div>
                            </div>
                        </div>
                    </div>

                    {/* <!-- Social & Copy Container --> */}
                    <div
                        className="flex flex-col items-start justify-between space-y-4 text-gray-250"
                    >
                        {/* <!-- icons Container --> */}
                        <div
                            className="flex items-center justify-center mx-auto space-x-4 md:justify-end md:mx-0"
                        >
                            {/* <!-- Icon 1 --> */}
                            <div className="h-8 group">
                                <a href="#">
                                <Image
                                    src="/icons/icon-facebook.svg"
                                    alt="icon-facebook"
                                    width={25}
                                    height={25}
                                    className="h-6"
                                ></Image>
                                </a>
                            </div>
                            {/* <!-- Icon 2 --> */}
                            <div className="h-8 group">
                                <a href="#">
                                <Image
                                    src="/icons/icon-twitter.svg"
                                    alt="icon-twitter"
                                    width={25}
                                    height={25}
                                    className="h-6"
                                ></Image>
                                </a>
                            </div>
                            {/* <!-- Icon 3 --> */}
                            <div className="h-8 group">
                                <a href="#">
                                <Image
                                    src="/icons/icon-pinterest.svg"
                                    alt="icon-pinterest"
                                    width={25}
                                    height={25}
                                    className="h-6"
                                ></Image>
                                
                                </a>
                            </div>
                            {/* <!-- Icon 4 --> */}
                            <div className="h-8 group">
                                <a href="#">
                                <Image
                                    src="/icons/icon-instagram.svg"
                                    alt="icon-instagram"
                                    width={25}
                                    height={25}
                                    className="h-6"
                                ></Image>
                                
                                </a>
                            </div>
                        </div>

                        {/* <!-- Copy --> */}
                        <div className="font-bold text-white">
                            &copy; 2023 UET VNU. All Rights Reserved
                        </div>
                    </div>
                </div>
            </div>
        </footer>
    )
}
