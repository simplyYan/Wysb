import os
import sys
import argparse
import importlib.util
from numba import jit, int32, int64, float32, float64, njit
import platform
import time
import random
import threading
import http
import configparser
import logging
import socket
import concurrent
import math
import json
import re
import string
import ctypes

class wysb:
    class resfilt:
        def __init__(self, func):
            self.func = func
            self.cache = {}

        def __call__(self, *args):
            if args in self.cache:
                return self.cache[args]
            result = self.func(*args)
            self.cache[args] = result
            return result

class TypeError(Exception):
    pass

class WysbCompiler:
    def __init__(self):
        self.variables = {}
        self.constants = {}
        self.imported_modules = {}

    def print(self, *args):
        print(*args)

    def check_type(self, value, expected_type):
        if expected_type == 'int32':
            if not isinstance(value, int):
                raise TypeError(f"Expected int32, got {type(value).__name__}")
            if not (-2**31 <= value < 2**31):
                raise TypeError("Value out of range for int32")
        elif expected_type == 'int64':
            if not isinstance(value, int):
                raise TypeError(f"Expected int64, got {type(value).__name__}")
            if not (-2**63 <= value < 2**63):
                raise TypeError("Value out of range for int64")
        elif expected_type == 'float32':
            if not isinstance(value, (float, int)):
                raise TypeError(f"Expected float32, got {type(value).__name__}")
        elif expected_type == 'float64':
            if not isinstance(value, (float, int)):
                raise TypeError(f"Expected float64, got {type(value).__name__}")
        elif expected_type == 'bool':
            if not isinstance(value, bool):
                raise TypeError(f"Expected bool, got {type(value).__name__}")
        elif expected_type == 'string':
            if not isinstance(value, str):
                raise TypeError(f"Expected string, got {type(value).__name__}")
        elif expected_type == 'uint':
            if not isinstance(value, int):
                raise TypeError(f"Expected uint, got {type(value).__name__}")
            if value < 0:
                raise TypeError("Value cannot be negative for uint")
        elif expected_type == 'char':
            if not isinstance(value, str) or len(value) != 1:
                raise TypeError("Expected single character for char")
        else:
            raise TypeError("Unknown type")

    def create_variable(self, var_name, var_value, var_type, mutable=True):
        self.check_type(var_value, var_type)
        if mutable:
            self.variables[var_name] = (var_value, var_type)
        else:
            self.constants[var_name] = (var_value, var_type)

    def get_variable(self, var_name):
        if var_name in self.constants:
            return self.constants[var_name]
        return self.variables.get(var_name, None)

    @njit
    def create_function(self, func_name, func_body):
        exec(f"def {func_name}(*args): {func_body}", globals())

    @njit
    def run_loop(self, loop_type, condition, body):
        exec(f"{loop_type} ({condition}): {body}", globals())

    def execute_conditional(self, condition, body):
        exec(f"if {condition}: {body}", globals())

    def execute_terminal(self, command):
        os.system(command)

    def get_platform(self):
        return platform.system()

    def measure_execution_time(self, func):
        start_time = time.time()
        func()
        end_time = time.time()
        return (end_time - start_time) * 1000  

    def arithmetic_operations(self, operation):
        return eval(operation)

    def allocate_memory(self, size):
        return bytearray(size)

    def create_array(self, *elements):
        return list(elements)

    def generate_random_number(self, lower, upper):
        return random.randint(lower, upper)

    def create_class(self, class_name, class_body):
        exec(f"class {class_name}: {class_body}", globals())

    def error_handling(self, try_body, except_body):
        try:
            exec(try_body, globals(), locals())
        except Exception as e:
            exec(f"{except_body}(e)", globals(), locals())

    def type_conversion(self, value, target_type):
        if target_type == 'int32':
            return int32(value)
        elif target_type == 'int64':
            return int64(value)
        elif target_type == 'float32':
            return float32(value)
        elif target_type == 'float64':
            return float64(value)
        elif target_type == 'bool':
            return bool(value)
        elif target_type == 'string':
            return str(value)
        elif target_type == 'uint':
            if value < 0:
                raise ValueError("Cannot convert negative value to uint")
            return value
        elif target_type == 'char':
            if len(value) != 1:
                raise ValueError("Cannot convert string of length > 1 to char")
            return value
        else:
            raise TypeError("Unknown target type")

    def file_crud(self, operation, filename, content=None):
        if operation == 'create':
            with open(filename, 'w') as file:
                file.write(content)
        elif operation == 'read':
            with open(filename, 'r') as file:
                return file.read()
        elif operation == 'update':
            with open(filename, 'a') as file:
                file.write(content)
        elif operation == 'delete':
            os.remove(filename)

    def import_module(self, module_name):
        self.imported_modules[module_name] = __import__(module_name)

    def import_wysb_module(self, filename):
        spec = importlib.util.spec_from_file_location("module", filename)
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        self.imported_modules[filename] = module

    def wysb_to_python(self, wysb_code):
        # Translate Wysb syntax to Python syntax
        wysb_code = wysb_code.replace("!func", "def")
        wysb_code = wysb_code.replace("!variable", "def")
        wysb_code = wysb_code.replace("<!--", "#")
        return wysb_code

    def execute_wysb_code(self, wysb_code):
        python_code = self.wysb_to_python(wysb_code)
        exec(python_code, globals())

    def compile_code(self, code):
        @jit(nopython=True)
        def compiled_code():
            exec(code)
        return compiled_code

    def convert_to_executable(self, wysb_file, exe_file, *modules):
        with open(wysb_file, 'r') as file:
            wysb_code = file.read()

        executable_code = f"""
import os
import sys
from numba import jit

def main():
    {self.wysb_to_python(wysb_code)}

if __name__ == "__main__":
    main()
"""

        temp_file = "temp_executable.py"
        with open(temp_file, 'w') as file:
            file.write(executable_code)

        os.system(f"pyinstaller --onefile --distpath . --name {exe_file} {temp_file}")

        os.remove(temp_file)
        os.remove(f"{exe_file}.spec")

    def execute_file(self, file_path):
        with open(file_path, 'r') as file:
            wysb_code = file.read()
        self.execute_wysb_code(wysb_code)

def main():
    parser = argparse.ArgumentParser(description='Wysb Compiler Tool')
    parser.add_argument('command', choices=['run', 'convert'], help='Command to execute or convert Wysb files.')
    parser.add_argument('file', help='The Wysb file to process.')
    parser.add_argument('--output', help='Output file name for conversion to executable.')

    args = parser.parse_args()

    compiler = WysbCompiler()

    if args.command == 'run':
        compiler.execute_file(args.file)

    elif args.command == 'convert':
        if not args.output:
            print("Output file name is required for conversion.")
            sys.exit(1)
        compiler.convert_to_executable(args.file, args.output)

if __name__ == "__main__":
    main()
