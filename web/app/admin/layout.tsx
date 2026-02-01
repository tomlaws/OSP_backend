import Link from 'next/link';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-white">
        {/* Top Navigation Grid */}
      <nav className="border-b-2 border-black">
        <div className="grid grid-cols-12 gap-0 h-16">
            {/* Logo Area */}
            <div className="col-span-3 border-r-2 border-black flex items-center px-6">
                <Link href="/admin" className="text-xl font-black tracking-tighter uppercase">
                  OSP <span className="text-[#D80000]">ADMIN</span>
                </Link>
            </div>
            
            {/* Navigation Links */}
            <div className="col-span-9 flex items-center justify-between px-6">
                <div className="flex space-x-8">
                    <Link
                    href="/admin"
                    className="text-sm font-bold uppercase tracking-wide hover:text-[#D80000] transition-colors"
                    >
                    Surveys
                    </Link>
                    <Link
                    href="/admin/create"
                    className="text-sm font-bold uppercase tracking-wide hover:text-[#D80000] transition-colors"
                    >
                    Create New
                    </Link>
                </div>
                <div className="flex items-center">
                    <span className="text-xs font-mono uppercase bg-black text-white px-2 py-1">Administrator_Mode</span>
                </div>
            </div>
        </div>
      </nav>
      
      {/* Main Content */}
      <main className="grid grid-cols-12 gap-0">
         {/* Sidebar / Grid line */}
         <div className="hidden md:block col-span-1 border-r-2 border-dashed border-gray-300 min-h-[calc(100vh-4rem)]"></div>
         
         <div className="col-span-12 md:col-span-11 p-8 md:p-12">
            {children}
         </div>
      </main>
    </div>
  );
}
