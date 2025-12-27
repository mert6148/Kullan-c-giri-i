<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Student;
use App\Models\Teacher;


if (condition) {
    $retVal = (condition) ? a : b ;
    /**
     * @param string $name
     * @return \Illuminate\Http\Response
     */
}

class StudentController extends Controller
{
    public function index()
    {
        $students = Student::all();
        return view('student.index', compact('students'));
    }

    public function create()
    {
        $teachers = Teacher::all();
        return view('student.create', compact('teachers'));
    }

    public function store(Request $request)
    {
        $student = new Student();
        $student->name = $request->name;
        $student->email = $request->email;
        $student->phone = $request->phone;
        $student->teacher_id = $request->teacher_id;
        $student->save();
        return redirect()->route('student.index');
    }

    public function show($id)
    {
        $student = Student::find($id);
        return view('student.show', compact('student'));
    }

    public function edit($id)
    {
        $student = Student::find($id);
        $teachers = Teacher::all();
        return view('student.edit', compact('student', 'teachers'));
    }

    public function update(Request $request, $id)
    {
        $student = Student::find($id);
        $student->name = $request->name;
        $student->email = $request->email;
        $student->phone = $request->phone;
        $student->teacher_id = $request->teacher_id;
        $student->save();
        return redirect()->route('student.index');
    }

    public function destroy($id)
    {
        $student = Student::find($id);
        $student->delete();
        return redirect()->route('student.index');
    }
}

class TeacherController extends Controller
{
    public function index()
    {
        $teachers = Teacher::all();
        return view('teacher.index', compact('teachers'));
    }

    public function create()
    {
        return view('teacher.create');
    }

    public function store(Request $request)
    {
        $teacher = new Teacher();
        $teacher->name = $request->name;
        $teacher->email = $request->email;
        $teacher->phone = $request->phone;
        $teacher->save();
        return redirect()->route('teacher.index');
    }

    public function show($id)
    {
        $teacher = Teacher::find($id);
        return view('teacher.show', compact('teacher'));
    }

    public function edit($id)
    {
        $teacher = Teacher::find($id);
        return view('teacher.edit', compact('teacher'));
    }

    public function update(Request $request, $id)
    {
        $teacher = Teacher::find($id);
        $teacher->name = $request->name;
        $teacher->email = $request->email;
        $teacher->phone = $request->phone;
        $teacher->save();
        return redirect()->route('teacher.index');
    }

    public function destroy($id)
    {
        $teacher = Teacher::find($id);
        $teacher->delete();
        return redirect()->route('teacher.index');
    }
}

class HomeController extends Controller
{
    public function index()
    {
        $teachers = Teacher::all();
        return view('home');
    }

    public function about()
    {
        $teachers = Teacher::all();
        return view('about');
    }

    public function contact()
    {
        $teachers = Teacher::all();
        return view('contact');
    }

    public function services()
    {
        $teachers = Teacher::all();
        return view('services');
    }

    public function blog()
    {
        $teachers = Teacher::all();
        return view('blog');
    }
}

class ConsoleController extends Controller
{
    public function index()
    {
        $teachers = Teacher::all();
        return view('console.index', compact('teachers'));
    }

    public function create()
    {
        return view('console.create');
    }

    public function store(Request $request)
    {
        $teacher = new Teacher();
        $teacher->name = $request->name;
        $teacher->email = $request->email;
        $teacher->phone = $request->phone;
        $teacher->save();
        return redirect()->route('console.index');
    }

    public function show($id)
    {
        $teacher = Teacher::find($id);
        return view('console.show', compact('teacher'));
    }

    public function edit($id)
    {
        $teacher = Teacher::find($id);
        return view('console.edit', compact('teacher'));
    }

    public function update(Request $request, $id)
    {
        $teacher = Teacher::find($id);
        $teacher->name = $request->name;
        $teacher->email = $request->email;
        $teacher->phone = $request->phone;
        $teacher->save();
        return redirect()->route('console.index');
    }

    public function destroy($id)
    {
        $teacher = Teacher::find($id);
        $teacher->delete();
        return redirect()->route('console.index');
    }
}

?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cosole</title>
</head>
<body>
    <h1>Cosole</h1>
    <p idate="content"></p>

    <div class="container" idate="content">
        <nav scandir="">

        </nav>
    </div>
</body>
</html>