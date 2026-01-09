# Example: Understanding Business Process

Complete walkthrough of investigating "How do I get an order to completed state for testing?"

## Initial Question

**User asks**: "I'm trying to test these edits in the dev1 environment. I've created a test order and am trying to cycle it through the process to the point where the green Dispute an Issue button appears. I think the order has to be completed, but I can't remember how to get it to that state. Help me figure this out."

**Research goal**: Understand order lifecycle, state transitions, and how to manually advance order to completed state in dev environment

## Step 1: Search for Process Documentation

Start with Glean to find documentation:

```
mcp__glean__company_search
query: "order lifecycle state machine completed"
datasources: ["confluence", "gdrive"]
```

**Results found**:
- Confluence: "Order State Machine Documentation"
- Confluence: "Testing Guide: Order States"
- Google Doc: "Fulfillment Workflows Overview"

## Step 2: Read State Machine Documentation

Get the comprehensive guide:

```bash
atl confluence search "order state machine"
```

**Found**: "Order State Machine Documentation" (ID: 98765432)

```bash
atl confluence page view 98765432 --output markdown
```

**Key information extracted**:

**Order states**:
1. `pending` - Order created, awaiting confirmation
2. `confirmed` - Payment processed, awaiting fulfillment
3. `in_progress` - Being prepared by caterer
4. `ready` - Ready for delivery/pickup
5. `in_transit` - Out for delivery
6. `completed` - Successfully delivered
7. `cancelled` - Order cancelled

**Dispute button appears**: After order reaches `completed` state

**State transitions**:
- `pending` → `confirmed`: Payment processing
- `confirmed` → `in_progress`: Caterer accepts order
- `in_progress` → `ready`: Caterer marks ready
- `ready` → `in_transit`: Delivery started
- `in_transit` → `completed`: Delivery confirmed

## Step 3: Find Testing Documentation

Check testing guide:

```bash
atl confluence search "testing guide order states"
```

**Found**: "Testing Guide: Order States" (ID: 98765433)

```bash
atl confluence page view 98765433 --output markdown
```

**Key information for dev environment**:

**Manual state advancement** (dev/staging only):
- Available in admin panel: `/admin/orders/:id/states`
- Can skip states for testing
- Must have `admin` role in dev

**Quick testing path**:
1. Create order
2. Admin panel → Orders → Find order
3. State Management → Advance to "completed"
4. Alternatively: Use Rails console commands

**Console commands provided**:
```ruby
# Find order
order = Order.find_by(number: "ORDER-NUMBER")

# Advance through states
order.confirm!
order.start_progress!
order.mark_ready!
order.start_transit!
order.complete!

# Or skip directly to completed
order.update!(state: 'completed', completed_at: Time.current)
```

## Step 4: Verify in Code

Confirm state machine implementation:

```bash
# Assume ez-rails checked out at ~/code/ezcater/ez/rails
cd ~/code/ezcater/ez/rails

# Find Order model
grep -r "state_machine" app/models/order.rb
```

View state machine definition:

```bash
cat app/models/order.rb | grep -A 50 "state_machine"
```

**Confirms**:
- States match documentation
- Methods `confirm!`, `start_progress!`, etc. exist
- `completed` state is valid

## Step 5: Find Dev Environment Access

Check how to access dev1 admin:

```bash
atl confluence search "dev1 admin access"
```

**Found**: "Development Environments Guide"

**Key info**:
- Dev1 URL: `https://dev1.sandbox-ezcater.net`
- Admin URL: `https://dev1.sandbox-ezcater.net/admin`
- Authentication: Use SSO with ezCater email
- Permissions: All engineers have admin role in dev

## Step 6: Check for Related Test Fixtures

See if there are automated ways:

```bash
cd ~/code/ezcater/ez/rails

# Find test factories
grep -r "completed" spec/factories/orders.rb
```

**Found factory**:
```ruby
factory :order do
  # ...

  trait :completed do
    state 'completed'
    completed_at { 1.day.ago }
  end
end
```

**Useful for**: Creating orders already in completed state for specs

## Step 7: Search for Recent Issues

Check if others had problems:

```bash
atl jira search "project = FX AND text ~ 'order completed state testing' AND created >= -90d"
```

**Found**: FX-4156 - "Dev environment orders stuck in pending"

```bash
atl jira issue view FX-4156
```

**Comments reveal**:
- Sometimes payments don't process in dev
- Workaround: Use console to force `confirmed` state first
- Payment processor stubs may need configuration

## Step 8: Search for Dispute Button Code

Understand when button appears:

```bash
cd ~/code/ezcater/ez/rails

# Find views mentioning dispute
grep -r "Dispute an Issue" app/views
```

**Found**: `app/views/orders/show.html.erb`

```ruby
<% if @order.completed? && @order.disputable? %>
  <%= button_to "Dispute an Issue", new_order_dispute_path(@order), class: "btn-success" %>
<% end %>
```

**Conditions**:
1. Order must be `completed?` (checks state == 'completed')
2. Order must be `disputable?` (check the method)

Check `disputable?` method:

```bash
grep -A 10 "def disputable?" app/models/order.rb
```

**Found**:
```ruby
def disputable?
  completed? &&
  completed_at > 30.days.ago &&
  !dispute.present?
end
```

**Additional requirements**:
- Completed within last 30 days
- No existing dispute

## Step 9: Check for Automated Testing Scripts

See if there are helper scripts:

```bash
cd ~/code/ezcater/ez/rails

find scripts -name "*order*" -o -name "*state*"
```

**Found**: `scripts/dev/advance_order_state.rb`

```bash
cat scripts/dev/advance_order_state.rb
```

**Script provides**:
```ruby
# Usage: rails runner scripts/dev/advance_order_state.rb ORDER-NUMBER completed
# Safely advances order to specified state with validation
```

Better than manual console commands!

## Step 10: Search Slack for Tips

Check team knowledge:

```
mcp__glean__company_search
query: "order completed state dev1 testing"
datasources: ["slack"]
```

**Found discussions** in #dev-support:
- Thread about payment stubs in dev1
- Recommendation to use `/admin/orders/:id/states` for manual testing
- Note: Console method `complete!` may trigger webhooks (use `update!` instead for testing)

## Synthesis: Complete Answer

### Immediate Solution

To get an order to completed state in dev1:

**Option 1: Admin Panel (Recommended)**

1. Go to `https://dev1.sandbox-ezcater.net/admin`
2. Navigate to Orders → Find your order
3. Click "State Management"
4. Select "completed" from dropdown
5. Click "Update State"
6. Verify `completed_at` is set (should be automatic)

**Option 2: Rails Console**

```bash
# SSH to dev1 pod or use eztilt
rails console

# Find and update order
order = Order.find_by(number: "YOUR-ORDER-NUMBER")
order.update!(state: 'completed', completed_at: Time.current)
```

**Option 3: Helper Script**

```bash
rails runner scripts/dev/advance_order_state.rb YOUR-ORDER-NUMBER completed
```

### Why Dispute Button Appears

Button shows when ALL conditions met:

1. ✅ Order state is `completed`
2. ✅ Completed within last 30 days (`completed_at > 30.days.ago`)
3. ✅ No existing dispute

### Complete Order Lifecycle

**Normal production flow**:
```
pending → confirmed → in_progress → ready → in_transit → completed
```

**Dev testing shortcut**:
```
pending → [skip to] completed
```

### State Transition Details

**What triggers each state** (production):

| From State | To State | Trigger | System Component |
|------------|----------|---------|------------------|
| pending | confirmed | Payment success | Payment processor |
| confirmed | in_progress | Caterer accepts | Caterer portal |
| in_progress | ready | Caterer marks ready | Caterer portal |
| ready | in_transit | Driver starts delivery | Driver app |
| in_transit | completed | Delivery confirmed | Driver app or customer |

**In dev environment**:
- Can skip states via admin panel
- Console commands also work
- Helper scripts provide safety checks

### Common Issues

**Order stuck in pending**:
- Cause: Payment processor stub not configured
- Fix: Use console to force `confirmed` first, then `completed`

**Dispute button not showing**:
- Check: `completed_at` must be set and recent (< 30 days)
- Check: No existing dispute on order
- Check: Order must be in `completed` state (not just `completed_at` present)

### Testing Best Practices

**From testing guide**:

1. **Use admin panel** for manual testing (safest)
2. **Avoid `complete!` method** in console (triggers webhooks)
3. **Use `update!` directly** for cleanest state change
4. **Check `disputable?`** after setting completed to verify button conditions
5. **Set `completed_at`** to recent timestamp (< 30 days)

## Final Report

### Executive Summary

To get an order to completed state in dev1, use the admin panel at `https://dev1.sandbox-ezcater.net/admin/orders/:id/states` and select "completed" from the State Management dropdown. The "Dispute an Issue" button will appear if the order was completed within the last 30 days and has no existing dispute.

### Step-by-Step Process

1. **Access admin panel**:
   - URL: `https://dev1.sandbox-ezcater.net/admin`
   - Login with SSO (ezCater email)

2. **Find order**:
   - Navigate to Orders
   - Search by order number

3. **Change state**:
   - Click "State Management"
   - Select "completed"
   - Click "Update State"
   - System automatically sets `completed_at`

4. **Verify dispute button**:
   - Return to order page
   - Green "Dispute an Issue" button should appear
   - If not visible, check conditions (see below)

### Alternative Methods

**Rails console** (if admin panel unavailable):
```ruby
order = Order.find_by(number: "ORDER-NUMBER")
order.update!(state: 'completed', completed_at: Time.current)
```

**Helper script** (safest console option):
```bash
rails runner scripts/dev/advance_order_state.rb ORDER-NUMBER completed
```

### Dispute Button Requirements

Button appears when:
- ✅ Order state == `completed`
- ✅ `completed_at` < 30 days ago
- ✅ No existing dispute

**Code location**: `app/views/orders/show.html.erb` calls `@order.disputable?`

**Check conditions**:
```ruby
order.completed?  # => true
order.completed_at > 30.days.ago  # => true
order.dispute.present?  # => false
```

### References

**Confluence**:
- "Order State Machine Documentation" (ID: 98765432) - Complete state definitions
- "Testing Guide: Order States" (ID: 98765433) - Dev environment procedures
- "Development Environments Guide" - Dev1 access info

**Code locations**:
- `app/models/order.rb` - State machine definition and `disputable?` method
- `app/views/orders/show.html.erb` - Dispute button conditional
- `scripts/dev/advance_order_state.rb` - Helper script for state changes

**Jira**:
- FX-4156 - Common dev environment issue with order states

**Glean**:
- Slack #dev-support - Team tips and troubleshooting

### Common Troubleshooting

**Q: Dispute button not showing after setting completed?**
1. Check `completed_at` is set: `order.completed_at`
2. Verify it's recent: `order.completed_at > 30.days.ago`
3. Check for existing dispute: `order.dispute.present?`

**Q: Order stuck in pending?**
Payment processor may not be configured. Force to confirmed first:
```ruby
order.update!(state: 'confirmed', confirmed_at: Time.current)
order.update!(state: 'completed', completed_at: Time.current)
```

**Q: Can I use `complete!` method?**
Not recommended in dev - triggers webhooks. Use `update!` instead.

## Research Techniques Demonstrated

1. **Documentation search**: Started with Confluence for official process docs
2. **Code verification**: Confirmed documentation matched actual implementation
3. **Test fixture discovery**: Found factory traits for automated testing
4. **Issue correlation**: Checked Jira for common problems
5. **Code search**: Used grep to find related views and conditional logic
6. **Script discovery**: Found helper scripts for common operations
7. **Team knowledge**: Leveraged Slack discussions for practical tips
8. **Condition analysis**: Traced button visibility logic through code

This example demonstrates how to understand a business process by combining documentation, code, issues, and team knowledge for a complete answer.
