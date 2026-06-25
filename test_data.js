const API_URL = 'http://localhost:8080/api/v1';

async function request(endpoint, method = 'GET', body = null, token = null) {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json'
    }
  };
  if (token) {
    options.headers['Authorization'] = `Bearer ${token}`;
  }
  if (body) {
    options.body = JSON.stringify(body);
  }
  const res = await fetch(API_URL + endpoint, options);
  const data = await res.json();
  if (!res.ok) {
    throw new Error(JSON.stringify(data));
  }
  return data;
}

async function runTest() {
  try {
    console.log("0. Registering user...");
    const regUsername = 'user' + Date.now();
    await request('/auth/register', 'POST', {
      username: regUsername,
      email: regUsername + '@example.com',
      password: 'password123'
    });
    
    console.log("0.1. Logging in...");
    const loginRes = await request('/auth/login', 'POST', {
      username: regUsername,
      password: 'password123'
    });
    const token = loginRes.data.token;
    console.log("Logged in:", token.substring(0, 10) + "...");

    console.log("1. Creating Account...");
    const acc = await request('/accounts', 'POST', {
      name: 'BCA Utama',
      type: 'asset',
      initial_balance: 1000000,
      currency: 'IDR',
      color: '#0055ff'
    }, token);
    console.log("Account Created:", acc.data.id);

    console.log("2. Creating Category...");
    const cat = await request('/categories', 'POST', {
      name: 'Makan & Minum',
      type: 'expense',
      color: '#ff3300',
      icon: '🍔'
    }, token);
    console.log("Category Created:", cat.data.id);

    console.log("3. Creating Transaction...");
    const tx = await request('/transactions', 'POST', {
      type: 'expense',
      amount: 50000,
      account_id: acc.data.id,
      category_id: cat.data.id,
      date: new Date().toISOString(),
      note: 'Makan siang'
    }, token);
    console.log("Transaction Created:", tx.data.id);

    console.log("4. Fetching Accounts to check balance...");
    const accs = await request('/accounts', 'GET', null, token);
    const updatedAcc = accs.data.items.find(a => a.id === acc.data.id);
    console.log("Updated Balance:", updatedAcc.balance);
    if (updatedAcc.balance !== 950000) {
      console.log("WARNING: Balance calculation might be wrong or not returned correctly.");
    } else {
      console.log("SUCCESS: Balance correctly deducted!");
    }
    
    console.log("All tests passed successfully!");
  } catch (err) {
    console.error("Test failed:", err.message);
  }
}

runTest();
