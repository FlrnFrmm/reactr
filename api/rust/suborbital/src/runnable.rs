pub mod default_runnable;

use std::mem;
use std::slice;
use crate::error::Error;
use crate::util;
use crate::STATE;

extern {
	fn return_result(result_pointer: *const u8, result_size: i32, ident: i32);
	fn return_error(code: i32, result_pointer: *const u8, result_size: i32, ident: i32);
}

pub trait Runnable {
	fn run(&self, input: Vec<u8>) -> Result<Vec<u8>, Error>;
}

pub fn use_runnable(runnable: &'static dyn Runnable) {
	unsafe {
		STATE.runnable = runnable;
	}
}

/// # Safety
/// 
/// We hand over the the pointer to the allocated memory.
/// Caller has to ensure that the memory gets freed again.
#[no_mangle]
pub unsafe extern fn allocate(size: i32) -> *const u8 {
	let mut buffer = Vec::with_capacity(size as usize);

	let pointer = buffer.as_mut_ptr();

	mem::forget(buffer);

	pointer as *const u8
}

/// # Safety
#[no_mangle]
pub unsafe extern fn deallocate(pointer: *const u8, size: i32) {
	let _ = slice::from_raw_parts(pointer, size as usize);
}

/// # Safety
#[no_mangle]
pub unsafe extern fn run_e(pointer: *const u8, size: i32, ident: i32) {
	STATE.ident = ident;

	// rebuild the memory into something usable
	let in_slice: &[u8] = slice::from_raw_parts(pointer, size as usize);

	let in_bytes = Vec::from(in_slice);

	let mut code = 0;

	// call the runnable and check its result
	let result: Vec<u8> = match STATE.runnable.run(in_bytes) {
		Ok(val) => val,
		Err(e) => {
			match e {
    			Error::Run { code: c, message } => {
					code = c;
					util::to_vec(message)
				},
    			Error::Host { message } => {
					code = -1;
					util::to_vec(message)
				},
			}
		}
	};

	let result_slice = result.as_slice();
	let result_size = result_slice.len();

	// call back to reactr to return the result or error
	if code != 0 {
		return_error(code, result_slice.as_ptr() as *const u8, result_size as i32, ident);
	} else {
		return_result(result_slice.as_ptr() as *const u8, result_size as i32, ident);
	}
}
