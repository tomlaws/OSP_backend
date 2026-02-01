import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen grid grid-cols-1 md:grid-cols-12 gap-0 relative overflow-hidden">
        {/* Abstract Geometric Element */}
        <div className="hidden md:block absolute top-0 right-0 w-1/3 h-full bg-[#D80000] -z-10"></div>
        
        {/* Left Column Content */}
        <div className="col-span-12 md:col-span-8 p-12 md:p-24 flex flex-col justify-center">
            
            {/* Header / Brand */}
            <div className="mb-24">
                <h2 className="text-xl font-bold tracking-tight uppercase border-b-2 border-black pb-4 inline-block">
                    OSP &mdash; 2026
                </h2>
            </div>

            {/* Main Title */}
            <h1 className="text-6xl md:text-8xl font-black tracking-tighter leading-none mb-12">
                SURVEY<br/>PLATFORM_
            </h1>

            {/* Subtext */}
            <div className="flex flex-col md:flex-row gap-12 text-lg md:w-3/4">
                <p className="font-medium leading-relaxed">
                    A systematic approach to data collection. 
                    Objective insights powered by artificial intelligence. 
                    Precision in every response.
                </p>
                
                {/* Visual Grid Element */}
                <div className="hidden md:block w-px bg-black h-24 self-center"></div>

                <div className="flex flex-col gap-4 min-w-[200px]">
                     <Link
                        href="/admin"
                        className="group flex items-center justify-between border-2 border-black px-6 py-4 font-bold uppercase hover:bg-black hover:text-white transition-colors duration-200"
                    >
                        <span>Dashboard</span>
                        <span className="group-hover:translate-x-2 transition-transform">&rarr;</span>
                    </Link>
                </div>
            </div>
        </div>

        {/* Right decoration (Mobile only) */}
        <div className="md:hidden h-24 bg-[#D80000] w-full"></div>
    </div>
  );
}
